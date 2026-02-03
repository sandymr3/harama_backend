DROP POLICY IF EXISTS tenant_exams_isolation ON exams;
DROP POLICY IF EXISTS tenant_submissions_isolation ON submissions;
DROP POLICY IF EXISTS tenant_feedback_isolation ON feedback_events;

ALTER TABLE exams DISABLE ROW LEVEL SECURITY;
ALTER TABLE submissions DISABLE ROW LEVEL SECURITY;
ALTER TABLE feedback_events DISABLE ROW LEVEL SECURITY;
