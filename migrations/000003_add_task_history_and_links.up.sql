CREATE TABLE IF NOT EXISTS task_events (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    field_name TEXT NOT NULL DEFAULT '',
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_task_events_task_id_created_at ON task_events (task_id, created_at DESC);

CREATE TABLE IF NOT EXISTS task_links (
    id BIGSERIAL PRIMARY KEY,
    source_task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    target_task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    link_type TEXT NOT NULL CHECK (link_type IN ('relates_to', 'blocks', 'duplicates')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT task_links_distinct_tasks CHECK (source_task_id <> target_task_id),
    CONSTRAINT task_links_unique UNIQUE (source_task_id, target_task_id, link_type)
);

CREATE INDEX IF NOT EXISTS idx_task_links_source_task_id ON task_links (source_task_id);
CREATE INDEX IF NOT EXISTS idx_task_links_target_task_id ON task_links (target_task_id);
