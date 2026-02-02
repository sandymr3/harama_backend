package service

import (
	"context"
	"fmt"
	"harama/internal/repository/postgres"
	"harama/internal/segmentation"
	"harama/internal/storage"

	"github.com/google/uuid"
)

type SegmentationService struct {
	repo      *postgres.SubmissionRepo
	detector  *segmentation.DiagramDetector
	storage   *storage.MinioStorage
}

func NewSegmentationService(repo *postgres.SubmissionRepo, detector *segmentation.DiagramDetector, storage *storage.MinioStorage) *SegmentationService {
	return &SegmentationService{
		repo:     repo,
		detector: detector,
		storage:  storage,
	}
}

func (s *SegmentationService) ExtractDiagrams(ctx context.Context, submissionID uuid.UUID, pageImage []byte) ([]string, error) {
	rects, err := s.detector.DetectRegions(pageImage)
	if err != nil {
		return nil, err
	}

	var diagramURLs []string
	for i, rect := range rects {
		cropped, err := s.detector.ExtractRegion(pageImage, rect)
		if err != nil {
			return nil, err
		}

		objectName := fmt.Sprintf("submissions/%s/diagram_%d.png", submissionID.String(), i)
		url, err := s.storage.UploadFile(ctx, objectName, cropped, "image/png")
		if err != nil {
			return nil, err
		}
		diagramURLs = append(diagramURLs, url)
	}

	return diagramURLs, nil
}
