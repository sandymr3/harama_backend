package ai

import (
    "context"
    "harama/internal/domain"
)

type Provider interface {
    Grade(ctx context.Context, req GradingRequest) (domain.GradingResult, error)
}

type GradingRequest struct {
    Answer      domain.AnswerSegment
    Rubric      domain.Rubric
    EvaluatorID string
}
