ALTER TABLE tasks
    ADD COLUMN due_date DATE,
    ADD COLUMN due_time TIME,
    ADD CONSTRAINT tasks_due_time_requires_date CHECK (due_date IS NOT NULL OR due_time IS NULL);

UPDATE tasks
SET due_date = due_at::date,
    due_time = due_at::time
WHERE due_at IS NOT NULL;

DROP INDEX IF EXISTS idx_tasks_due_at;
ALTER TABLE tasks DROP COLUMN due_at;
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks (due_date);
