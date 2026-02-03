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
	examRepo   *postgres.ExamRepo
	aiProvider ai.Provider
}

func NewFeedbackService(repo *postgres.FeedbackRepo, gradeRepo *postgres.GradeRepo, examRepo *postgres.ExamRepo, aiProvider ai.Provider) *FeedbackService {
	return &FeedbackService{
		repo:       repo,
		gradeRepo:  gradeRepo,
		examRepo:   examRepo,
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

func (s *FeedbackService) AnalyzeQuestionPatterns(ctx context.Context, questionID uuid.UUID) (ai.AnalysisResult, error) {
	// 1. Get question and rubric
	question, err := s.examRepo.GetQuestionByID(ctx, questionID)
	if err != nil {
		return ai.AnalysisResult{}, err
	}

	if question.Rubric == nil {
		return ai.AnalysisResult{}, nil // No rubric to analyze against
	}

	// 2. Get feedback events for this question
	events, err := s.repo.GetFeedbackByQuestion(ctx, questionID)
	if err != nil {
		return ai.AnalysisResult{}, err
	}

	if len(events) == 0 {
		return ai.AnalysisResult{}, nil // Not enough data
	}

	// 3. Call AI to analyze patterns
	return s.aiProvider.AnalyzePatterns(ctx, ai.AnalysisRequest{
		QuestionID: questionID,
		Rubric:     *question.Rubric,
		Events:     events,
	})
}

func (s *FeedbackService) AdaptRubric(ctx context.Context, questionID uuid.UUID) error {
	// 1. Analyze patterns
	analysis, err := s.AnalyzeQuestionPatterns(ctx, questionID)
	if err != nil {
		return err
	}

	if analysis.Recommendation == "" {
		return nil
	}

	// 2. Get the current question and rubric
	question, err := s.examRepo.GetQuestionByID(ctx, questionID)
	if err != nil {
		return err
	}

	if question.Rubric == nil {
		return nil
	}

	// In a real system, we'd call Gemini again to "Apply these recommendations to this JSON rubric"
	// For this phase, we'll mark it as a task to be refined in Phase 5 corrections.
	
	// Example: Log the recommendation for now or update a 'GradingNotes' field
	question.Rubric.GradingNotes = "AI Recommendation: " + analysis.Recommendation

	return s.examRepo.UpdateRubric(ctx, question.Rubric)
}

func (s *FeedbackService) GetFeedbackByQuestion(ctx context.Context, questionID uuid.UUID) ([]domain.FeedbackEvent, error) {
	return s.repo.GetFeedbackByQuestion(ctx, questionID)
}