package jobs

import (
	"context"
	"harama/internal/service"
	"github.com/google/uuid"
)

type OCRJob struct {
	SubmissionID uuid.UUID
	Service      *service.OCRService
}

func (j *OCRJob) Execute(ctx context.Context) error {
	return j.Service.ProcessSubmission(ctx, j.SubmissionID)
}

func (j *OCRJob) ID() string {
	return "ocr-" + j.SubmissionID.String()
}

type GradingJob struct {
	SubmissionID uuid.UUID
	Service      *service.GradingService
}

func (j *GradingJob) Execute(ctx context.Context) error {
	return j.Service.GradeSubmission(ctx, j.SubmissionID)
}

func (j *GradingJob) ID() string {
	return "grading-" + j.SubmissionID.String()
}
