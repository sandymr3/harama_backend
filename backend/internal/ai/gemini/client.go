package gemini

import (
    "context"
    "fmt"
    
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
    
    "harama/internal/ai"
    "harama/internal/domain"
)

type Client struct {
    client *genai.Client
    model  *genai.GenerativeModel
}

func NewClient(apiKey string) (*Client, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }
    
    model := client.GenerativeModel("gemini-3-pro")
    model.SetTemperature(0.2)
    model.SetTopK(40)
    model.SetTopP(0.95)
    
    return &Client{
        client: client,
        model:  model,
    }, nil
}

func (c *Client) Grade(ctx context.Context, req ai.GradingRequest) (domain.GradingResult, error) {
	// Load appropriate prompt template
	promptTemplate := loadPromptTemplate(req.EvaluatorID)

	// Build prompt
	prompt := buildGradingPrompt(promptTemplate, req.Answer, req.Rubric)

	// Prepare parts for multimodal input
	parts := []genai.Part{genai.Text(prompt)}

	// Add diagrams if present (Phase 2)
	for _, diagramURL := range req.Answer.Diagrams {
		// In a real scenario, we might need to download the image from MinIO
		// or pass the bytes directly if we have them.
		// For now, we'll assume the request might include bytes or we use a helper.
		// Mock: adding a placeholder if we were to pass image data
		// parts = append(parts, genai.ImageData("image/png", diagramBytes))
		_ = diagramURL // Use URL to fetch or placeholder
	}

	// Call Gemini
	resp, err := c.model.GenerateContent(ctx, parts...)
	if err != nil {
		return domain.GradingResult{}, fmt.Errorf("gemini API error: %w", err)
	}

	// Parse structured response
	var result domain.GradingResult
	if err := parseResponse(resp, &result); err != nil {
		return domain.GradingResult{}, err
	}

	result.AIEvaluatorID = req.EvaluatorID
	return result, nil
}

func (c *Client) GenerateFeedback(ctx context.Context, req ai.FeedbackRequest) (string, error) {
	prompt := fmt.Sprintf(`
ROLE: Educational feedback specialist

STUDENT: %s
CURRENT GRADE: %.1f/%.1f
AI REASONING: %s

HISTORY OF OVERRIDES:
%v

TASK:
Generate personalized, encouraging, and actionable feedback for the student.
Focus on:
1. What they did well.
2. Specific areas for improvement based on current mistakes and history.
3. Actionable next steps.

TONE: Encouraging but honest.
LENGTH: 3-4 sentences.
`, req.StudentName, req.Grade.FinalScore, req.Grade.Confidence, req.Grade.Reasoning, req.History)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		return string(text), nil
	}

	return "", fmt.Errorf("unexpected response part type")
}

func (c *Client) AnalyzePatterns(ctx context.Context, req ai.AnalysisRequest) (ai.AnalysisResult, error) {
	prompt := fmt.Sprintf(`
ROLE: Educational pattern analyst

QUESTION ID: %s
RUBRIC: %v
FEEDBACK EVENTS:
%v

TASK:
Analyze the feedback events where teachers have overridden AI scores.
Identify:
1. Systematic biases (e.g., AI is consistently too strict on units).
2. Common themes in teacher corrections.
3. Recommendations for improving the rubric or system prompt.

OUTPUT JSON FORMAT:
{
  "patterns": ["pattern 1", "pattern 2"],
  "common_reasons": ["reason 1", "reason 2"],
  "recommendation": "detailed recommendation"
}
`, req.QuestionID, req.Rubric, req.Events)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return ai.AnalysisResult{}, fmt.Errorf("gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ai.AnalysisResult{}, fmt.Errorf("empty response from Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)
	if !ok {
		return ai.AnalysisResult{}, fmt.Errorf("unexpected response part type")
	}

	var result ai.AnalysisResult
	if err := json.Unmarshal([]byte(textPart), &result); err != nil {
		// If it's not JSON, we might need a more robust parser or just return as recommendation
		return ai.AnalysisResult{
			Recommendation: string(textPart),
		}, nil
	}

	return result, nil
}
