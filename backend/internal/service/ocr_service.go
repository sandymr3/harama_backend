package service

import (
	"context"
	"fmt"
	"harama/internal/domain"
	"harama/internal/repository/postgres"
	"harama/internal/storage"

	"github.com/google/uuid"
)

// OCRProcessor defines the contract for text extraction strategies
type OCRProcessor interface {
	ExtractText(ctx context.Context, fileBytes []byte, mimeType string) (*domain.OCRResult, error)
}

type OCRService struct {
	repo    *postgres.SubmissionRepo
	storage *storage.MinioStorage
	processor OCRProcessor
}

func NewOCRService(repo *postgres.SubmissionRepo, storage *storage.MinioStorage, processor OCRProcessor) *OCRService {
	return &OCRService{
		repo:      repo,
		storage:   storage,
		processor: processor,
	}
}

func (s *OCRService) ProcessSubmission(ctx context.Context, submissionID uuid.UUID) error {
	// 1. Get submission metadata
	sub, err := s.repo.GetByID(ctx, submissionID)
	if err != nil {
		return err
	}

	// 2. Mark as processing
	err = s.repo.UpdateStatus(ctx, submissionID, domain.StatusProcessing)
	if err != nil {
		return err
	}

	// 3. Process each OCR result
	var finalResults []domain.OCRResult
	for _, res := range sub.OCRResults {
		// Assume ImageURL is the object name in MinIO
		imgBytes, err := s.storage.GetFile(ctx, res.ImageURL)
		if err != nil {
			return fmt.Errorf("failed to get file from storage: %w", err)
		}

		// TODO: Determine mime type from file extension
		mimeType := "image/png" 
		ocrResult, err := s.processor.ExtractText(ctx, imgBytes, mimeType)
		if err != nil {
			return fmt.Errorf("failed to extract text: %w", err)
		}

		ocrResult.PageNumber = res.PageNumber
		ocrResult.ImageURL = res.ImageURL
		finalResults = append(finalResults, *ocrResult)
	}

	return s.repo.SaveOCRResults(ctx, submissionID, finalResults)
}
