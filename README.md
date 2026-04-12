# todolist-api

Todo list backend in Go with PostgreSQL, REST CRUD endpoints, SQL migrations,
structured request logging, health probes, Swagger docs, and container support
for a React frontend.

## Stack

- Go HTTP API with `chi`
- PostgreSQL via `pgx`
- SQL migrations via `golang-migrate` in `cmd/migrate`
- Structured JSON logging with `log/slog`
- Swagger/OpenAPI 2.0 support with `swag` annotations and `http-swagger`
- Docker multi-stage build

## Project Structure

```text
.
├── cmd
│   ├── api
│   └── migrate
├── docs
├── internal
│   ├── config
│   ├── db
│   ├── http
│   ├── response
│   └── tasks
├── migrations
├── Dockerfile
├── Makefile
└── docker-compose.yml
```

## Configuration

Copy `.env.example` values into your shell or your preferred env file loader.

| Variable | Description | Default |
| --- | --- | --- |
| `APP_ADDR` | HTTP bind address | `:8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/todolist?sslmode=disable` |
| `ALLOWED_ORIGIN` | CORS origin for the React frontend | `*` |
| `READINESS_TIMEOUT` | DB ping timeout for readiness checks | `2s` |

## Local Development

Start PostgreSQL:

```bash
docker compose up -d postgres
```

Apply schema migrations:

```bash
make migrate-up
```

Run the API:

```bash
make run
```

The API will be available on `http://localhost:8080`.

## Data Model

Tasks now support richer workflow metadata:

- `status`: `todo`, `in_progress`, `done`, `archived`
- `priority`: `low`, `medium`, `high`
- `due_at`: optional due date/time
- `completed_at`: automatically set when a task moves to `done`
- `parent_task_id`: optional parent reference for sub-tasks

The longer task details remain in `description`.

## Database Migrations

Migrations live in [`migrations`](./migrations).

Commands:

```bash
make migrate-up
make migrate-down
make migrate-reset
make migrate-version
```

These commands run the Go migration utility in [`cmd/migrate`](./cmd/migrate),
so you do not need to install a separate migration CLI.

### Current Schema

The task table starts in [`migrations/000001_create_tasks.up.sql`](./migrations/000001_create_tasks.up.sql)
and is extended in [`migrations/000002_extend_tasks.up.sql`](./migrations/000002_extend_tasks.up.sql).
After all migrations, the table shape is effectively:

```sql
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'done', 'archived')),
    priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    due_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    parent_task_id BIGINT REFERENCES tasks(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## API Endpoints

Base path: `/api/v1`

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/v1/tasks` | List all tasks |
| `GET` | `/api/v1/tasks/{taskID}` | Fetch one task |
| `POST` | `/api/v1/tasks` | Create a task |
| `PUT` | `/api/v1/tasks/{taskID}` | Update a task |
| `DELETE` | `/api/v1/tasks/{taskID}` | Delete a task |
| `GET` | `/health/live` | Liveness probe |
| `GET` | `/health/ready` | Readiness probe backed by DB ping |

### Example Requests

Create a parent task:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Launch the frontend",
    "description": "Prepare the React app for production",
    "status": "todo",
    "priority": "high",
    "due_at": "2026-04-20T09:00:00Z"
  }'
```

Create a sub-task:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hook React to the API",
    "description": "Wire the task screens to the backend endpoints",
    "status": "in_progress",
    "priority": "medium",
    "parent_task_id": 1
  }'
```

List tasks:

```bash
curl http://localhost:8080/api/v1/tasks
```

Get a task by ID:

```bash
curl http://localhost:8080/api/v1/tasks/1
```

Update a task:

```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Launch the frontend",
    "description": "React app is wired to the backend and ready to deploy",
    "status": "done",
    "priority": "high",
    "due_at": "2026-04-20T09:00:00Z",
    "parent_task_id": null
  }'
```

Delete a task:

```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```

## Logging

Every request is logged with structured JSON fields including method, path,
status, duration, and remote address. That makes it straightforward to ship
logs to a collector in Docker, Kubernetes, or a cloud runtime.

## Swagger Documentation

Swagger annotations are in the Go handlers and the repo includes a checked-in
`docs` package so the Swagger UI can be served immediately.

Swagger UI endpoint:

```text
http://localhost:8080/swagger/index.html
```

The raw spec is available at:

```text
http://localhost:8080/swagger/doc.json
```

To regenerate the Swagger files after editing handlers or annotations:

1. Install the generator:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Regenerate:

```bash
make swagger
```

## Docker

Build the image:

```bash
docker build -t todolist-api .
```

Run the container:

```bash
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/todolist?sslmode=disable" \
  -e ALLOWED_ORIGIN="http://localhost:3000" \
  todolist-api
```

On container startup, the entrypoint runs `todolist-migrate up` before starting
the API process. If there are no pending migrations, startup continues normally.
If a migration fails, the container exits instead of starting the server with an
unexpected schema.

To skip startup migrations for a container instance:

```bash
docker run --rm -p 8080:8080 \
  -e RUN_MIGRATIONS=false \
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/todolist?sslmode=disable" \
  todolist-api
```

### Health Probes

For orchestrators such as Kubernetes, use:

- Liveness probe: `GET /health/live`
- Readiness probe: `GET /health/ready`

This pattern is fine for local Docker and single-instance deployments. In a
multi-replica environment, a separate migration job or init container is often
cleaner.

## Notes

- `status=done` automatically sets `completed_at`.
- Moving a task away from `done` clears `completed_at`.
- `parent_task_id` must reference an existing task and cannot reference the
  task itself.
- Readiness returns `503` if PostgreSQL is unavailable.
- `ALLOWED_ORIGIN` should be tightened to your real frontend origin outside
  local development.
