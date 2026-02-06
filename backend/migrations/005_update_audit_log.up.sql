ALTER TABLE audit_log ADD COLUMN actor_type VARCHAR(20);
ALTER TABLE audit_log ADD COLUMN metadata JSONB;
ALTER TABLE audit_log ADD COLUMN hash VARCHAR(64);

ALTER TABLE grades ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Populate hash for existing rows if any (using a dummy value or actual hash if preferred)
-- Since it's probably empty or has few rows, we can just set a dummy one or leave it null if we change the column to not null later.
-- For now, let's keep it nullable or provide a default.
UPDATE audit_log SET hash = 'initial' WHERE hash IS NULL;

-- If we want to enforce it later
-- ALTER TABLE audit_log ALTER COLUMN hash SET NOT NULL;
