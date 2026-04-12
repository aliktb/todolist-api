ALTER TABLE tasks ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE tasks
SET completed = (status = 'done');

DROP INDEX IF EXISTS idx_tasks_parent_task_id;
DROP INDEX IF EXISTS idx_tasks_due_at;
DROP INDEX IF EXISTS idx_tasks_status;

ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS tasks_priority_check,
    DROP CONSTRAINT IF EXISTS tasks_status_check,
    DROP COLUMN parent_task_id,
    DROP COLUMN completed_at,
    DROP COLUMN due_at,
    DROP COLUMN priority,
    DROP COLUMN status;

CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks (completed);
