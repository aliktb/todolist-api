ALTER TABLE tasks
    ADD COLUMN status TEXT NOT NULL DEFAULT 'todo',
    ADD COLUMN priority TEXT NOT NULL DEFAULT 'medium',
    ADD COLUMN due_at TIMESTAMPTZ,
    ADD COLUMN completed_at TIMESTAMPTZ,
    ADD COLUMN parent_task_id BIGINT REFERENCES tasks(id) ON DELETE SET NULL,
    ADD CONSTRAINT tasks_status_check CHECK (status IN ('todo', 'in_progress', 'done', 'archived')),
    ADD CONSTRAINT tasks_priority_check CHECK (priority IN ('low', 'medium', 'high'));

UPDATE tasks
SET status = CASE
        WHEN completed THEN 'done'
        ELSE 'todo'
    END,
    completed_at = CASE
        WHEN completed THEN updated_at
        ELSE NULL
    END;

DROP INDEX IF EXISTS idx_tasks_completed;
ALTER TABLE tasks DROP COLUMN completed;

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_at ON tasks (due_at);
CREATE INDEX IF NOT EXISTS idx_tasks_parent_task_id ON tasks (parent_task_id);
