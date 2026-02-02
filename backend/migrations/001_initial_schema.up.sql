CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE exams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    title VARCHAR(255) NOT NULL,
    subject VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exam_id UUID NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    points INTEGER NOT NULL,
    answer_type VARCHAR(50) NOT NULL,
    visual_aids JSONB
);

CREATE TABLE rubrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    full_credit_criteria JSONB NOT NULL,
    partial_credit_rules JSONB NOT NULL,
    common_mistakes JSONB NOT NULL,
    key_concepts JSONB,
    grading_notes TEXT,
    strict_mode BOOLEAN DEFAULT FALSE
);

CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exam_id UUID NOT NULL REFERENCES exams(id),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processing_status VARCHAR(50) NOT NULL,
    ocr_results JSONB,
    answers JSONB
);

CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL REFERENCES submissions(id),
    question_id UUID NOT NULL REFERENCES questions(id),
    score DECIMAL(5,2),
    max_score INTEGER NOT NULL,
    confidence DECIMAL(3,2),
    reasoning TEXT,
    criteria_met JSONB,
    mistakes_found JSONB,
    ai_evaluator_id VARCHAR(100),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    actor_id UUID, 
    changes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
