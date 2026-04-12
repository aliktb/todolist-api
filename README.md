# todolist-api

> [!INFO]
> This is a vibe-coded project.
>
> That means:
>
> - Some parts are genius
> - Some parts are... mysterious 🌀
> - Tests may or may not exist
> - Refactoring is a future concept

Todo list backend in Go with PostgreSQL, REST CRUD endpoints, SQL migrations,
structured request logging, history tracking, task linking, Swagger docs, and
container support for a React frontend: [todolist-ui](https://github.com/aliktb/todolist-ui)

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

Tasks support workflow metadata and optional relationships:

- `status`: `todo`, `in_progress`, `done`, `archived`
- `priority`: `low`, `medium`, `high`
- `due_date`: optional due date
- `due_time`: optional time on the selected due date; cannot be set without `due_date`
- `completed_at`: automatically set when a task moves to `done`
- `parent_task_id`: optional parent reference for sub-tasks

There are also two supporting tables:

- `task_events`: append-only history for task changes and link activity
- `task_links`: explicit links between related tasks

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

### Schema Overview

The task schema is created in
[`migrations/000001_create_tasks.up.sql`](./migrations/000001_create_tasks.up.sql),
extended in
[`migrations/000002_extend_tasks.up.sql`](./migrations/000002_extend_tasks.up.sql),
and the history/link tables are added in
[`migrations/000003_add_task_history_and_links.up.sql`](./migrations/000003_add_task_history_and_links.up.sql),
with due dates and optional due times introduced in
[`migrations/000004_split_due_date_and_time.up.sql`](./migrations/000004_split_due_date_and_time.up.sql).

Key tables after all migrations:

```sql
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'todo',
    priority TEXT NOT NULL DEFAULT 'medium',
    due_date DATE,
    due_time TIME,
    CHECK (due_date IS NOT NULL OR due_time IS NULL),
    completed_at TIMESTAMPTZ,
    parent_task_id BIGINT REFERENCES tasks(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE task_events (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    field_name TEXT NOT NULL DEFAULT '',
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE task_links (
    id BIGSERIAL PRIMARY KEY,
    source_task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    target_task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    link_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## API Endpoints

Base path: `/api/v1`

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/v1/tasks` | List all tasks |
| `GET` | `/api/v1/tasks/search?q=...` | Fuzzy-find tasks for linking |
| `GET` | `/api/v1/tasks/{taskID}` | Fetch one task |
| `POST` | `/api/v1/tasks` | Create a task |
| `PUT` | `/api/v1/tasks/{taskID}` | Update a task |
| `DELETE` | `/api/v1/tasks/{taskID}` | Delete a task |
| `GET` | `/api/v1/tasks/{taskID}/history` | Fetch task change history |
| `GET` | `/api/v1/tasks/{taskID}/links` | List links for a task |
| `POST` | `/api/v1/tasks/{taskID}/links` | Create a task link |
| `DELETE` | `/api/v1/tasks/{taskID}/links/{linkID}` | Delete a task link |
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
    "due_date": "2026-04-20",
    "due_time": "09:00"
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

Search tasks for linking:

```bash
curl "http://localhost:8080/api/v1/tasks/search?q=frontend"
```

Update a task and trigger history entries:

```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Launch the frontend",
    "description": "React app is wired to the backend and ready to deploy",
    "status": "done",
    "priority": "high",
    "due_date": "2026-04-20",
    "due_time": null,
    "parent_task_id": null
  }'
```

View task history:

```bash
curl http://localhost:8080/api/v1/tasks/1/history
```

Create a task link:

```bash
curl -X POST http://localhost:8080/api/v1/tasks/1/links \
  -H "Content-Type: application/json" \
  -d '{
    "target_task_id": 2,
    "link_type": "blocks"
  }'
```

List links for a task:

```bash
curl http://localhost:8080/api/v1/tasks/1/links
```

Delete a task link:

```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1/links/1
```

## History Model

Task history is exposed via `GET /api/v1/tasks/{taskID}/history` and is backed by
an append-only `task_events` table. The API records:

- task creation via `task_created`
- field changes via `task_updated`
- task link creation via `task_linked`
- task link deletion via `task_unlinked`

Each history row includes `old_value` and `new_value` as JSON so your frontend
can render a timeline or diff view.

## Linking Model

Task links are separate from `parent_task_id` so you can support both sub-tasks
and peer-to-peer relationships.

Supported `link_type` values:

- `relates_to`
- `blocks`
- `duplicates`

The frontend can use `/api/v1/tasks/search?q=...` to find candidate tasks and
then call `POST /api/v1/tasks/{taskID}/links` to create the actual link.

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

```bash
go install github.com/swaggo/swag/cmd/swag@latest
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
- `due_time` cannot be provided unless `due_date` is also provided.
- Moving a task away from `done` clears `completed_at`.
- `parent_task_id` must reference an existing task and cannot reference the
  task itself.
- Task links require a valid target task and cannot self-reference.
- Readiness returns `503` if PostgreSQL is unavailable.
- `ALLOWED_ORIGIN` should be tightened to your real frontend origin outside
  local development.
