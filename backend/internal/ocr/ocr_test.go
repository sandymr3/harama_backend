package ocr

import (
	"context"
	"os"
	"testing"
)

func TestGeminiOCRWithAnswerJpeg(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	processor, err := NewGeminiOCRProcessor(apiKey)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}
	defer processor.Close()

	imgBytes, err := os.ReadFile("../../answer.jpeg")
	if err != nil {
		t.Fatalf("Failed to read answer.jpeg: %v", err)
	}

	result, err := processor.ExtractText(context.Background(), imgBytes, "image/jpeg")
	if err != nil {
		t.Fatalf("OCR failed: %v", err)
	}
	t.Logf("Extracted: %s", result.RawText)
}
