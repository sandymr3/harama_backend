package service

import (
	"context"
	"harama/internal/ai"
	"harama/internal/domain"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
)

type FeedbackService struct {
	repo       *postgres.FeedbackRepo
	gradeRepo  *postgres.GradeRepo
	aiProvider ai.Provider
}

func NewFeedbackService(repo *postgres.FeedbackRepo, gradeRepo *postgres.GradeRepo, aiProvider ai.Provider) *FeedbackService {
	return &FeedbackService{
		repo:       repo,
		gradeRepo:  gradeRepo,
		aiProvider: aiProvider,
	}
}

func (s *FeedbackService) CaptureOverrideFeedback(ctx context.Context, submissionID uuid.UUID, questionID uuid.UUID, teacherScore float64, teacherReason string) error {
	// 1. Get the existing grade to find the AI score and reasoning
	grades, err := s.gradeRepo.GetBySubmission(ctx, submissionID)
	if err != nil {
		return err
	}

	var originalGrade *domain.FinalGrade
	for _, g := range grades {
		if g.QuestionID == questionID {
			originalGrade = &g
			break
		}
	}

	if originalGrade == nil {
		return nil // Or return error if grade should exist
	}

	aiScore := 0.0
	if originalGrade.AIScore != nil {
		aiScore = *originalGrade.AIScore
	}

	event := &domain.FeedbackEvent{
		ID:            uuid.New(),
		QuestionID:    questionID,
		SubmissionID:  submissionID,
		AIScore:       aiScore,
		TeacherScore:  teacherScore,
		Delta:         teacherScore - aiScore,
		AIReasoning:   originalGrade.Reasoning,
		TeacherReason: teacherReason,
	}

	return s.repo.SaveFeedbackEvent(ctx, event)
}

func (s *FeedbackService) GenerateStudentFeedback(ctx context.Context, submissionID uuid.UUID, questionID uuid.UUID, studentName string) (string, error) {
	// 1. Get current grade
	grades, err := s.gradeRepo.GetBySubmission(ctx, submissionID)
	if err != nil {
		return "", err
	}

	var currentGrade *domain.FinalGrade
	for _, g := range grades {
		if g.QuestionID == questionID {
			currentGrade = &g
			break
		}
	}

	if currentGrade == nil {
		return "", nil // Grade not found
	}

	// 2. Get historical feedback for this student (simplified to same submission for now)
	history, err := s.repo.GetFeedbackBySubmission(ctx, submissionID)
	if err != nil {
		return "", err
	}

	// 3. Call AI to generate feedback
	return s.aiProvider.GenerateFeedback(ctx, ai.FeedbackRequest{
		Grade:       *currentGrade,
		History:     history,
		StudentName: studentName,
	})
}

func (s *FeedbackService) GetFeedbackByQuestion(ctx context.Context, questionID uuid.UUID) ([]domain.FeedbackEvent, error) {
	return s.repo.GetFeedbackByQuestion(ctx, questionID)
}
