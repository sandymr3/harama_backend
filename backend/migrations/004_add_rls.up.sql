-- Enable Row Level Security
ALTER TABLE exams ENABLE ROW LEVEL SECURITY;
ALTER TABLE submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE feedback_events ENABLE ROW LEVEL SECURITY;

-- Note: In a real app, we'd also have 'tenants' and 'users' tables
-- and use current_setting('app.current_tenant') set by the app.

-- Create policies (assuming tenant_id exists on all core tables)
-- For exams:
CREATE POLICY tenant_exams_isolation ON exams
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::UUID);

-- For submissions:
-- If submissions doesn't have tenant_id directly, we might need to join, 
-- but for performance, we usually add tenant_id to all core tables.
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS tenant_id UUID;
CREATE POLICY tenant_submissions_isolation ON submissions
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::UUID);

-- For feedback_events:
ALTER TABLE feedback_events ADD COLUMN IF NOT EXISTS tenant_id UUID;
CREATE POLICY tenant_feedback_isolation ON feedback_events
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::UUID);
