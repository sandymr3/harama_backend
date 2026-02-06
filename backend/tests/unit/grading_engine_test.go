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

func (m *MockProvider) RefineRubric(ctx context.Context, req ai.RefineRubricRequest) (domain.Rubric, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(domain.Rubric), args.Error(1)
}

func TestEngine_MultiEvaluatorGrade(t *testing.T) {
	mockAI := new(MockProvider)
	engine := grading.NewEngine(mockAI)

	ctx := context.Background()
	answer := domain.AnswerSegment{Text: "The mitochondria is the powerhouse of the cell"}
	
	rubric := domain.Rubric{
		FullCreditCriteria: []domain.Criterion{
			{ID: "calc_correct", Points: 4.0},
			{ID: "method_correct", Points: 3.0},
			{ID: "units_correct", Points: 1.0},
			{ID: "explanation_clear", Points: 2.0},
		},
	}

	// Scenario 1: High Consensus
	// All evaluators find all criteria met -> Score 10
	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{
		Score:       10, // AI suggests 10
		MaxScore:    10,
		Confidence:  0.95,
		CriteriaMet: []string{"calc_correct", "method_correct", "units_correct", "explanation_clear"},
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
	// Rubric Enforcer (Strict): Only explanation_clear -> Score 2
	mockAI.On("Grade", ctx, mock.MatchedBy(func(req ai.GradingRequest) bool {
		return req.EvaluatorID == "rubric_enforcer"
	})).Return(domain.GradingResult{
		Score:       2,
		MaxScore:    10,
		Confidence:  0.9,
		CriteriaMet: []string{"explanation_clear"},
	}, nil)

	// Reasoning Validator (Lenient): calc + method + explanation -> Score 9
	mockAI.On("Grade", ctx, mock.MatchedBy(func(req ai.GradingRequest) bool {
		return req.EvaluatorID == "reasoning_validator"
	})).Return(domain.GradingResult{
		Score:       9,
		MaxScore:    10,
		Confidence:  0.8,
		CriteriaMet: []string{"calc_correct", "method_correct", "explanation_clear"},
	}, nil)

	// Structural Analyzer: method + explanation -> Score 5
	mockAI.On("Grade", ctx, mock.MatchedBy(func(req ai.GradingRequest) bool {
		return req.EvaluatorID == "structural_analyzer"
	})).Return(domain.GradingResult{
		Score:       5,
		MaxScore:    10,
		Confidence:  0.85,
		CriteriaMet: []string{"method_correct", "explanation_clear"},
	}, nil)

	finalGrade, multiEval, err = engine.GradeAnswer(ctx, answer, rubric, "Science", "Question")

	assert.NoError(t, err)
	assert.True(t, multiEval.ShouldEscalate, "Should escalate due to high variance between 2, 9, 5")
	assert.Equal(t, domain.GradeStatusReview, finalGrade.Status)

	// Scenario 3: Moderate Variance (Consensus reached)
	mockAI = new(MockProvider)
	engine = grading.NewEngine(mockAI)

	// 1: calc + method -> 7
	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{
		Score:       7, 
		MaxScore:    10, 
		Confidence:  0.9,
		CriteriaMet: []string{"calc_correct", "method_correct"},
	}, nil).Once()

	// 2: calc + method + units -> 8
	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{
		Score:       8, 
		MaxScore:    10, 
		Confidence:  0.9,
		CriteriaMet: []string{"calc_correct", "method_correct", "units_correct"},
	}, nil).Once()

	// 3: calc + method -> 7
	mockAI.On("Grade", ctx, mock.Anything).Return(domain.GradingResult{
		Score:       7, 
		MaxScore:    10, 
		Confidence:  0.9,
		CriteriaMet: []string{"calc_correct", "method_correct"},
	}, nil).Once()

	finalGrade, multiEval, err = engine.GradeAnswer(ctx, answer, rubric, "Science", "Question")

	assert.NoError(t, err)
	assert.False(t, multiEval.ShouldEscalate, "Variance should be low enough")
	assert.Contains(t, multiEval.Reasoning, "High confidence in consensus") 
	
	// Mean is 7.33, expect result around there
	assert.InDelta(t, 7.33, finalGrade.FinalScore, 0.1)
}
