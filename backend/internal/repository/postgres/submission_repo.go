package postgres

import (
	"context"
	"harama/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubmissionRepo struct {
	db *bun.DB
}

func NewSubmissionRepo(db *bun.DB) *SubmissionRepo {
	return &SubmissionRepo{db: db}
}

func (r *SubmissionRepo) Create(ctx context.Context, sub *domain.Submission) error {
	_, err := r.db.NewInsert().Model(sub).Exec(ctx)
	return err
}

func (r *SubmissionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Submission, error) {
	sub := new(domain.Submission)
	err := r.db.NewSelect().
		Model(sub).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (r *SubmissionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProcessingStatus) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Submission)(nil)).
		Set("processing_status = ?", status).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *SubmissionRepo) SaveOCRResults(ctx context.Context, id uuid.UUID, results []domain.OCRResult) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Submission)(nil)).
		Set("ocr_results = ?", results).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
