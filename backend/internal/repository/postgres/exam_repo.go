package postgres

import (
	"context"
	"harama/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ExamRepo struct {
	db *bun.DB
}

func NewExamRepo(db *bun.DB) *ExamRepo {
	return &ExamRepo{db: db}
}

func (r *ExamRepo) Create(ctx context.Context, exam *domain.Exam) error {
	_, err := r.db.NewInsert().Model(exam).Exec(ctx)
	return err
}

func (r *ExamRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Exam, error) {
	exam := new(domain.Exam)
	err := r.db.NewSelect().
		Model(exam).
		Relation("Questions").
		Relation("Questions.Rubric").
		Where("e.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return exam, nil
}

func (r *ExamRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := r.db.NewSelect().
		Model(&exams).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Scan(ctx)
	return exams, err
}

func (r *ExamRepo) CreateQuestion(ctx context.Context, question *domain.Question) error {
	_, err := r.db.NewInsert().Model(question).Exec(ctx)
	return err
}

func (r *ExamRepo) UpdateRubric(ctx context.Context, rubric *domain.Rubric) error {
	_, err := r.db.NewInsert().
		Model(rubric).
		On("CONFLICT (question_id) DO UPDATE").
		Set("full_credit_criteria = EXCLUDED.full_credit_criteria").
		Set("partial_credit_rules = EXCLUDED.partial_credit_rules").
		Set("common_mistakes = EXCLUDED.common_mistakes").
		Set("key_concepts = EXCLUDED.key_concepts").
		Set("grading_notes = EXCLUDED.grading_notes").
		Set("strict_mode = EXCLUDED.strict_mode").
		Exec(ctx)
	return err
}
