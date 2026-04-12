package tasks

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("task not found")
	ErrInvalidParent = errors.New("invalid parent task")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context) ([]Task, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, status, priority, due_at, completed_at, parent_task_id, created_at, updated_at
		FROM tasks
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueAt,
			&task.CompletedAt,
			&task.ParentTaskID,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int64) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, status, priority, due_at, completed_at, parent_task_id, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueAt,
		&task.CompletedAt,
		&task.ParentTaskID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	return task, err
}

func (r *Repository) Create(ctx context.Context, input CreateTaskInput) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, `
		INSERT INTO tasks (title, description, status, priority, due_at, parent_task_id, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, CASE WHEN $3 = 'done' THEN NOW() ELSE NULL END)
		RETURNING id, title, description, status, priority, due_at, completed_at, parent_task_id, created_at, updated_at
	`, input.Title, input.Description, input.Status, input.Priority, input.DueAt, input.ParentTaskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueAt,
		&task.CompletedAt,
		&task.ParentTaskID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return Task{}, translateError(err)
	}
	return task, nil
}

func (r *Repository) Update(ctx context.Context, id int64, input UpdateTaskInput) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, `
		UPDATE tasks
		SET title = $2,
		    description = $3,
		    status = $4,
		    priority = $5,
		    due_at = $6,
		    parent_task_id = $7,
		    completed_at = CASE
		        WHEN $4 = 'done' THEN COALESCE(completed_at, NOW())
		        ELSE NULL
		    END,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, status, priority, due_at, completed_at, parent_task_id, created_at, updated_at
	`, id, input.Title, input.Description, input.Status, input.Priority, input.DueAt, input.ParentTaskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueAt,
		&task.CompletedAt,
		&task.ParentTaskID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, translateError(err)
	}
	return task, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	commandTag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func translateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrInvalidParent
	}
	return err
}
