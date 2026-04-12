# AGENTS.md

## Purpose

This repository contains a Go backend for a todo list API backed by PostgreSQL.
It exposes REST endpoints for task CRUD, history, search, and task links, and it
serves Swagger docs for frontend/API integration.

## Stack

- Go
- `chi` for routing
- PostgreSQL via `pgx`
- SQL migrations via `cmd/migrate`
- Swagger docs in `docs/`
- Docker multi-stage image build

## Important Paths

- `cmd/api`: API entrypoint
- `cmd/migrate`: migration runner
- `internal/http`: HTTP handlers and middleware
- `internal/tasks`: task domain models and repository logic
- `migrations`: SQL migrations
- `docs`: checked-in Swagger artifacts
- `README.md`: setup and API documentation

## Common Commands

Run the API:

```bash
make run
```

Run migrations:

```bash
make migrate-up
make migrate-down
make migrate-version
```

Run tests:

```bash
env GOCACHE=$(pwd)/.cache/go-build GOPATH=$(pwd)/.cache/gopath go test ./...
```

Regenerate Swagger files:

```bash
make swagger
```

## Development Notes

- Keep migrations additive and reversible where practical.
- If you change request/response models or routes, update Swagger and `README.md`.
- Task history is stored in `task_events`; changes to task mutation logic should
  preserve event recording inside the same transaction.
- Task links are separate from `parent_task_id`; do not conflate sub-task
  hierarchy with peer task relationships.
- Prefer small, atomic commits with imperative, capitalized commit messages.

## API Notes

Current task-related capabilities include:

- CRUD for tasks
- task search for fuzzy linking
- per-task history endpoint
- per-task link CRUD endpoints

When adding new task fields, consider whether they should also:

- appear in search results
- be tracked in history events
- be documented in Swagger and README examples
