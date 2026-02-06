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

func (r *GradeRepo) SaveEscalation(ctx context.Context, escalation *domain.EscalationCase) error {
	_, err := r.db.NewInsert().Model(escalation).Exec(ctx)
	return err
}

type QuestionStat struct {
	QuestionID    uuid.UUID `bun:"question_id"`
	AvgScore      float64   `bun:"avg_score"`
	ScoreVariance float64   `bun:"score_variance"`
	ZeroScores    int       `bun:"zero_scores"`
	PerfectScores int       `bun:"perfect_scores"`
	TotalGraded   int       `bun:"total_graded"`
}

func (r *GradeRepo) GetExamStats(ctx context.Context, examID uuid.UUID) ([]QuestionStat, error) {
	var stats []QuestionStat
	// Join grades with questions to filter by exam_id
	// Note: We need to handle max_score comparison carefully. 
	// Assuming max_score is in grades table (it is per schema)
	
	err := r.db.NewSelect().
		Table("grades").
		ColumnExpr("grades.question_id").
		ColumnExpr("AVG(grades.final_score) as avg_score").
		ColumnExpr("STDDEV(grades.final_score) as score_variance").
		ColumnExpr("COUNT(CASE WHEN grades.final_score = 0 THEN 1 END) as zero_scores").
		ColumnExpr("COUNT(CASE WHEN grades.final_score = grades.max_score THEN 1 END) as perfect_scores").
		ColumnExpr("COUNT(*) as total_graded").
		Join("JOIN questions ON grades.question_id = questions.id").
		Where("questions.exam_id = ?", examID).
		Group("grades.question_id").
		Scan(ctx, &stats)
		
	return stats, err
}
