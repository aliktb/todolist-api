DROP INDEX IF EXISTS idx_task_links_target_task_id;
DROP INDEX IF EXISTS idx_task_links_source_task_id;
DROP TABLE IF EXISTS task_links;

DROP INDEX IF EXISTS idx_task_events_task_id_created_at;
DROP TABLE IF EXISTS task_events;
