package unit_test

import (
	"context"
	"testing"

	"harama/internal/ai"
	"harama/internal/domain"
	"harama/internal/grading"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProvider is a testify mock for ai.Provider
type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) Grade(ctx context.Context, req ai.GradingRequest) (domain.GradingResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(domain.GradingResult), args.Error(1)
}

func (m *MockProvider) GenerateFeedback(ctx context.Context, req ai.FeedbackRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockProvider) AnalyzePatterns(ctx context.Context, req ai.AnalysisRequest) (ai.AnalysisResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(ai.AnalysisResult), args.Error(1)
}

func TestEngine_MultiEvaluatorGrade(t *testing.T) {
	mockAI := new(MockProvider)
	engine := grading.NewEngine(mockAI)

	ctx := context.Background()
	answer := domain.AnswerSegment{Text: "The mitochondria is the powerhouse of the cell"}
	rubric := domain.Rubric{}

	// Scenario 1: High Consensus
	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{
		Score:      10,
		MaxScore:   10,
		Confidence: 0.95,
	}, nil).Times(3)

	finalGrade, multiEval, err := engine.GradeAnswer(ctx, answer, rubric, "Science", "What is the mitochondria?")

	assert.NoError(t, err)
	assert.InDelta(t, 10.0, finalGrade.FinalScore, 0.001)
	assert.False(t, multiEval.ShouldEscalate)
	assert.Greater(t, multiEval.Confidence, 0.9)
	
	mockAI.AssertExpectations(t)
	
	// Reset mock for next scenario
	mockAI = new(MockProvider)
	engine = grading.NewEngine(mockAI)

	// Scenario 2: High Variance (Should Escalate)
	// Rubric Enforcer (Strict)
	mockAI.On("Grade", ctx, mock.MatchedBy(func(req ai.GradingRequest) bool {
		return req.EvaluatorID == "rubric_enforcer"
	})).Return(domain.GradingResult{Score: 2, MaxScore: 10, Confidence: 0.9}, nil)

	// Reasoning Validator (Lenient)
	mockAI.On("Grade", ctx, mock.MatchedBy(func(req ai.GradingRequest) bool {
		return req.EvaluatorID == "reasoning_validator"
	})).Return(domain.GradingResult{Score: 9, MaxScore: 10, Confidence: 0.8}, nil)

	// Structural Analyzer
	mockAI.On("Grade", ctx, mock.MatchedBy(func(req ai.GradingRequest) bool {
		return req.EvaluatorID == "structural_analyzer"
	})).Return(domain.GradingResult{Score: 5, MaxScore: 10, Confidence: 0.85}, nil)

	finalGrade, multiEval, err = engine.GradeAnswer(ctx, answer, rubric, "Science", "Question")

	assert.NoError(t, err)
	assert.True(t, multiEval.ShouldEscalate, "Should escalate due to high variance between 2 and 9")
	assert.Equal(t, domain.GradeStatusReview, finalGrade.Status)

	// Scenario 3: Moderate Variance (Consensus reached)
	mockAI = new(MockProvider)
	engine = grading.NewEngine(mockAI)

	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{Score: 6, MaxScore: 10, Confidence: 0.9}, nil).Once()
	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{Score: 9, MaxScore: 10, Confidence: 0.9}, nil).Once()
	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{Score: 7.5, MaxScore: 10, Confidence: 0.9}, nil).Once()

	finalGrade, multiEval, err = engine.GradeAnswer(ctx, answer, rubric, "Science", "Question")

	assert.NoError(t, err)
	assert.False(t, multiEval.ShouldEscalate, "Variance 1.5 should not trigger escalation (threshold > 1.5)")
	assert.Contains(t, multiEval.Reasoning, "Moderate variance")
	assert.InDelta(t, 7.5, finalGrade.FinalScore, 0.1)
}
