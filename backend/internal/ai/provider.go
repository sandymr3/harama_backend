package ai

import (
    "context"
    "harama/internal/domain"
)

type Provider interface {
    Grade(ctx context.Context, req GradingRequest) (domain.GradingResult, error)
    GenerateFeedback(ctx context.Context, req FeedbackRequest) (string, error)
    AnalyzePatterns(ctx context.Context, req AnalysisRequest) (AnalysisResult, error)
}

type GradingRequest struct {
    Answer      domain.AnswerSegment
    Rubric      domain.Rubric
    EvaluatorID string
}

type FeedbackRequest struct {
    Grade    domain.FinalGrade
    History  []domain.FeedbackEvent
    StudentName string
}

type AnalysisRequest struct {
    QuestionID uuid.UUID
    Rubric     domain.Rubric
    Events     []domain.FeedbackEvent
}

type AnalysisResult struct {
    Patterns       []string `json:"patterns"`
    CommonReasons  []string `json:"common_reasons"`
    Recommendation string   `json:"recommendation"`
}
