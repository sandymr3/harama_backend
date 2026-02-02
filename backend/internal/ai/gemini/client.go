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
    
    // Call Gemini
    resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
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

// Stubs for helper functions
func loadPromptTemplate(evaluatorID string) string { return "" }
func buildGradingPrompt(template string, answer domain.AnswerSegment, rubric domain.Rubric) string { return "" }
func parseResponse(resp *genai.GenerateContentResponse, result *domain.GradingResult) error { return nil }
