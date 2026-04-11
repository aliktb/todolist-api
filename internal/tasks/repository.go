package tasks

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("task not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context) ([]Task, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
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
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int64) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`, id).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	return task, err
}

func (r *Repository) Create(ctx context.Context, input CreateTaskInput) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description, completed, created_at, updated_at
	`, input.Title, input.Description).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}

func (r *Repository) Update(ctx context.Context, id int64, input UpdateTaskInput) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, `
		UPDATE tasks
		SET title = $2,
		    description = $3,
		    completed = $4,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, completed, created_at, updated_at
	`, id, input.Title, input.Description, input.Completed).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	return task, err
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
