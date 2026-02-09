package ocr

import (
	"context"
	"fmt"
	"strings"

	"harama/internal/domain"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiOCRProcessor struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiOCRProcessor(apiKey string) (*GeminiOCRProcessor, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	
	// Use Gemini 3 Flash Preview for faster/cheaper OCR
	model := client.GenerativeModel("gemini-3-flash-preview") 
	model.SetTemperature(0.1) // Low temperature for deterministic transcription
	
	return &GeminiOCRProcessor{
		client: client,
		model:  model,
	}, nil
}

func (p *GeminiOCRProcessor) ExtractText(ctx context.Context, fileBytes []byte, mimeType string) (*domain.OCRResult, error) {
	// Construct multimodal prompt
	prompt := genai.Text("Transcribe the handwritten or printed text in this exam page exactly as it appears. Do not correct spelling. Return only the transcribed text.")
	
	var imgData genai.Part
	
	// Default to image/png if mimeType is empty
	if mimeType == "" {
		mimeType = "image/png"
	}
	
	imgData = genai.Blob{
		MIMEType: mimeType,
		Data:     fileBytes,
	}

	resp, err := p.model.GenerateContent(ctx, prompt, imgData)
	if err != nil {
		return nil, fmt.Errorf("gemini ocr error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini OCR")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from Gemini OCR")
	}

	return &domain.OCRResult{
		RawText:    strings.TrimSpace(string(text)),
		Confidence: 0.90, // Gemini doesn't give token-level confidence easily in standard response, defaulting
	}, nil
}

func (p *GeminiOCRProcessor) Close() error {
	return p.client.Close()
}
