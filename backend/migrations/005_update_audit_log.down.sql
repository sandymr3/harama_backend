ALTER TABLE audit_log DROP COLUMN actor_type;
ALTER TABLE audit_log DROP COLUMN metadata;
ALTER TABLE audit_log DROP COLUMN hash;

ALTER TABLE grades DROP COLUMN updated_at;
