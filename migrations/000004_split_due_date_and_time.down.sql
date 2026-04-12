ALTER TABLE tasks ADD COLUMN due_at TIMESTAMPTZ;

UPDATE tasks
SET due_at = CASE
    WHEN due_date IS NULL THEN NULL
    WHEN due_time IS NULL THEN due_date::timestamp AT TIME ZONE 'UTC'
    ELSE (due_date::timestamp + due_time) AT TIME ZONE 'UTC'
END;

DROP INDEX IF EXISTS idx_tasks_due_date;
ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS tasks_due_time_requires_date,
    DROP COLUMN due_time,
    DROP COLUMN due_date;

CREATE INDEX IF NOT EXISTS idx_tasks_due_at ON tasks (due_at);
