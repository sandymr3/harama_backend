-- Add description to exams
ALTER TABLE exams ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '';

-- Add question_number and question_group to questions
ALTER TABLE questions ADD COLUMN IF NOT EXISTS question_number VARCHAR(50) DEFAULT '';
ALTER TABLE questions ADD COLUMN IF NOT EXISTS question_group VARCHAR(100) DEFAULT '';
