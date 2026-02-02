package postgres

import (
	"context"
	"harama/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GradeRepo struct {
	db *bun.DB
}

func NewGradeRepo(db *bun.DB) *GradeRepo {
	return &GradeRepo{db: db}
}

func (r *GradeRepo) SaveFinalGrade(ctx context.Context, grade *domain.FinalGrade) error {
	_, err := r.db.NewInsert().
		Model(grade).
		On("CONFLICT (submission_id, question_id) DO UPDATE").
		Set("final_score = EXCLUDED.final_score").
		Set("ai_score = EXCLUDED.ai_score").
		Set("override_score = EXCLUDED.override_score").
		Set("confidence = EXCLUDED.confidence").
		Set("status = EXCLUDED.status").
		Set("graded_by = EXCLUDED.graded_by").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func (r *GradeRepo) GetBySubmission(ctx context.Context, submissionID uuid.UUID) ([]domain.FinalGrade, error) {
	var grades []domain.FinalGrade
	err := r.db.NewSelect().
		Model(&grades).
		Where("submission_id = ?", submissionID).
		Scan(ctx)
	return grades, err
}

func (r *GradeRepo) CreateAuditLog(ctx context.Context, logEntry map[string]interface{}) error {
	_, err := r.db.NewInsert().
		Table("audit_log").
		Model(&logEntry).
		Exec(ctx)
	return err
}
