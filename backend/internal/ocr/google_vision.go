package ocr

import (
	"bytes"
	"context"
	"fmt"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
	"harama/internal/domain"
)

type GoogleVisionProcessor struct {
	client *vision.ImageAnnotatorClient
}

func NewGoogleVisionProcessor(apiKey string) (*GoogleVisionProcessor, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GoogleVisionProcessor{client: client}, nil
}

func (p *GoogleVisionProcessor) ExtractText(ctx context.Context, fileBytes []byte, mimeType string) (*domain.OCRResult, error) {
	// 1. Handle PDF files using Vision API's file annotation (Sync if small, usually Async for large)
	// For this synchronous grading flow, we will assume standard image formats are preferred,
	// BUT we will add the logic to handle PDFs by treating them as "files" if the API supports sync.
	// Note: Vision API "DetectDocumentText" handles images. For PDFs, we need "AnnotateFile".
	
	if mimeType == "application/pdf" {
		return p.extractPDFText(ctx, fileBytes)
	}

	// 2. Handle Image files (JPEG, PNG, etc.)
	image, err := vision.NewImageFromReader(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	annotation, err := p.client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return nil, err
	}

	if annotation == nil {
		return &domain.OCRResult{}, nil
	}

	return &domain.OCRResult{
		RawText:    annotation.Text,
		Confidence: 0.95,
	}, nil
}

func (p *GoogleVisionProcessor) extractPDFText(ctx context.Context, data []byte) (*domain.OCRResult, error) {
    // PDF handling in Google Vision usually requires GCS storage for async processing.
    // For a direct synchronous API, it's limited.
    // We will return a specific error prompting the user to upload images, 
    // OR (if we had `pdfcpu`) we would split it here.
    //
    // DECISION: For this prototype, we will return a "Not Implemented" error for PDFs
    // to reflect that we need the async GCS flow, avoiding a half-baked implementation.
    return nil, fmt.Errorf("direct PDF processing requires async GCS pipeline; please convert to images")
}

func (p *GoogleVisionProcessor) Close() error {
	return p.client.Close()
}
