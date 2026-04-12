package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("task not found")
	ErrInvalidParent = errors.New("invalid parent task")
	ErrLinkExists    = errors.New("task link already exists")
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
		if err := scanTask(rows, &task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *Repository) Search(ctx context.Context, query string, limit int) ([]TaskSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, status, priority, due_at
		FROM tasks
		WHERE title ILIKE '%' || $1 || '%'
		   OR description ILIKE '%' || $1 || '%'
		ORDER BY
			CASE WHEN title ILIKE $1 || '%' THEN 0 ELSE 1 END,
			POSITION(LOWER($1) IN LOWER(title)),
			id ASC
		LIMIT $2
	`, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskSummary
	for rows.Next() {
		var task TaskSummary
		if err := rows.Scan(&task.ID, &task.Title, &task.Status, &task.Priority, &task.DueAt); err != nil {
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

func (r *Repository) History(ctx context.Context, taskID int64) ([]TaskEvent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, task_id, event_type, field_name, old_value, new_value, created_at
		FROM task_events
		WHERE task_id = $1
		ORDER BY created_at DESC, id DESC
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []TaskEvent
	for rows.Next() {
		var event TaskEvent
		var oldValue []byte
		var newValue []byte
		if err := rows.Scan(&event.ID, &event.TaskID, &event.EventType, &event.FieldName, &oldValue, &newValue, &event.CreatedAt); err != nil {
			return nil, err
		}
		event.OldValue = json.RawMessage(oldValue)
		event.NewValue = json.RawMessage(newValue)
		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *Repository) Create(ctx context.Context, input CreateTaskInput) (Task, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Task{}, err
	}
	defer tx.Rollback(ctx)

	var task Task
	err = tx.QueryRow(ctx, `
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

	payload, err := marshalJSON(task)
	if err != nil {
		return Task{}, err
	}
	if err := insertTaskEvent(ctx, tx, task.ID, "task_created", "", nil, payload); err != nil {
		return Task{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *Repository) Update(ctx context.Context, id int64, input UpdateTaskInput) (Task, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Task{}, err
	}
	defer tx.Rollback(ctx)

	before, err := getTaskTx(ctx, tx, id)
	if err != nil {
		return Task{}, err
	}

	var task Task
	err = tx.QueryRow(ctx, `
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
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, translateError(err)
	}

	for _, event := range diffTaskEvents(before, task) {
		if err := insertTaskEvent(ctx, tx, task.ID, event.EventType, event.FieldName, event.OldValue, event.NewValue); err != nil {
			return Task{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Task{}, err
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

func (r *Repository) ListLinks(ctx context.Context, taskID int64) ([]TaskLink, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.id, l.source_task_id, l.target_task_id, l.link_type, l.created_at,
		       t.id, t.title, t.status, t.priority, t.due_at
		FROM task_links l
		JOIN tasks t ON t.id = CASE WHEN l.source_task_id = $1 THEN l.target_task_id ELSE l.source_task_id END
		WHERE l.source_task_id = $1 OR l.target_task_id = $1
		ORDER BY l.created_at DESC, l.id DESC
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []TaskLink
	for rows.Next() {
		var link TaskLink
		if err := rows.Scan(
			&link.ID,
			&link.SourceTaskID,
			&link.TargetTaskID,
			&link.LinkType,
			&link.CreatedAt,
			&link.LinkedTask.ID,
			&link.LinkedTask.Title,
			&link.LinkedTask.Status,
			&link.LinkedTask.Priority,
			&link.LinkedTask.DueAt,
		); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, rows.Err()
}

func (r *Repository) CreateLink(ctx context.Context, sourceTaskID int64, input CreateTaskLinkInput) (TaskLink, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return TaskLink{}, err
	}
	defer tx.Rollback(ctx)

	var link TaskLink
	err = tx.QueryRow(ctx, `
		INSERT INTO task_links (source_task_id, target_task_id, link_type)
		VALUES ($1, $2, $3)
		RETURNING id, source_task_id, target_task_id, link_type, created_at
	`, sourceTaskID, input.TargetTaskID, input.LinkType).Scan(
		&link.ID,
		&link.SourceTaskID,
		&link.TargetTaskID,
		&link.LinkType,
		&link.CreatedAt,
	)
	if err != nil {
		return TaskLink{}, translateError(err)
	}

	link.LinkedTask, err = getTaskSummaryTx(ctx, tx, input.TargetTaskID)
	if err != nil {
		return TaskLink{}, err
	}

	payload, err := marshalJSON(map[string]any{
		"link_id":        link.ID,
		"source_task_id": link.SourceTaskID,
		"target_task_id": link.TargetTaskID,
		"link_type":      link.LinkType,
	})
	if err != nil {
		return TaskLink{}, err
	}
	if err := insertTaskEvent(ctx, tx, sourceTaskID, "task_linked", "link", nil, payload); err != nil {
		return TaskLink{}, err
	}
	if sourceTaskID != input.TargetTaskID {
		if err := insertTaskEvent(ctx, tx, input.TargetTaskID, "task_linked", "link", nil, payload); err != nil {
			return TaskLink{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return TaskLink{}, err
	}
	return link, nil
}

func (r *Repository) DeleteLink(ctx context.Context, taskID, linkID int64) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sourceTaskID int64
	var targetTaskID int64
	var linkType string
	err = tx.QueryRow(ctx, `
		SELECT source_task_id, target_task_id, link_type
		FROM task_links
		WHERE id = $1 AND (source_task_id = $2 OR target_task_id = $2)
	`, linkID, taskID).Scan(&sourceTaskID, &targetTaskID, &linkType)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	payload, err := marshalJSON(map[string]any{
		"link_id":        linkID,
		"source_task_id": sourceTaskID,
		"target_task_id": targetTaskID,
		"link_type":      linkType,
	})
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM task_links WHERE id = $1`, linkID); err != nil {
		return err
	}

	if err := insertTaskEvent(ctx, tx, sourceTaskID, "task_unlinked", "link", payload, nil); err != nil {
		return err
	}
	if sourceTaskID != targetTaskID {
		if err := insertTaskEvent(ctx, tx, targetTaskID, "task_unlinked", "link", payload, nil); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func getTaskTx(ctx context.Context, tx pgx.Tx, id int64) (Task, error) {
	var task Task
	err := tx.QueryRow(ctx, `
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

func getTaskSummaryTx(ctx context.Context, tx pgx.Tx, id int64) (TaskSummary, error) {
	var task TaskSummary
	err := tx.QueryRow(ctx, `
		SELECT id, title, status, priority, due_at
		FROM tasks
		WHERE id = $1
	`, id).Scan(&task.ID, &task.Title, &task.Status, &task.Priority, &task.DueAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return TaskSummary{}, ErrNotFound
	}
	return task, err
}

func diffTaskEvents(before, after Task) []TaskEvent {
	var events []TaskEvent
	appendIfChanged := func(field string, oldValue any, newValue any) {
		oldPayload, err := marshalJSON(oldValue)
		if err != nil {
			return
		}
		newPayload, err := marshalJSON(newValue)
		if err != nil {
			return
		}
		if string(oldPayload) == string(newPayload) {
			return
		}
		events = append(events, TaskEvent{
			EventType: "task_updated",
			FieldName: field,
			OldValue:  oldPayload,
			NewValue:  newPayload,
		})
	}

	appendIfChanged("title", before.Title, after.Title)
	appendIfChanged("description", before.Description, after.Description)
	appendIfChanged("status", before.Status, after.Status)
	appendIfChanged("priority", before.Priority, after.Priority)
	appendIfChanged("due_at", before.DueAt, after.DueAt)
	appendIfChanged("completed_at", before.CompletedAt, after.CompletedAt)
	appendIfChanged("parent_task_id", before.ParentTaskID, after.ParentTaskID)
	return events
}

func insertTaskEvent(ctx context.Context, tx pgx.Tx, taskID int64, eventType, fieldName string, oldValue, newValue json.RawMessage) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO task_events (task_id, event_type, field_name, old_value, new_value)
		VALUES ($1, $2, $3, $4::jsonb, $5::jsonb)
	`, taskID, eventType, fieldName, nullableJSON(oldValue), nullableJSON(newValue))
	return err
}

func nullableJSON(value json.RawMessage) any {
	if len(value) == 0 {
		return nil
	}
	return string(value)
}

func marshalJSON(value any) (json.RawMessage, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(payload), nil
}

func scanTask(rows pgx.Rows, task *Task) error {
	return rows.Scan(
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
}

func translateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return ErrInvalidParent
		case "23505":
			return ErrLinkExists
		default:
			return fmt.Errorf("database error %s: %w", pgErr.Code, err)
		}
	}
	return err
}
