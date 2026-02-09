package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"harama/internal/config"
	"harama/internal/ocr"
	"harama/internal/repository/postgres"
	"harama/internal/service"
	"harama/internal/storage"

	"github.com/google/uuid"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Connect to DB
	db, err := postgres.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	// Get tenant
	var tenant struct {
		ID   uuid.UUID
		Name string
	}
	if err := db.NewSelect().Table("tenants").Column("id", "name").Limit(1).Scan(ctx, &tenant); err != nil {
		log.Fatal("No tenant found:", err)
	}
	log.Printf("Using Tenant: %s (%s)", tenant.Name, tenant.ID)

	// Initialize services
	minioClient, err := storage.NewMinioStorage(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, cfg.MinioUseSSL)
	if err != nil {
		log.Fatal("MinIO init failed:", err)
	}

	ocrProcessor, err := ocr.NewGeminiOCRProcessor(cfg.GeminiAPIKey)
	if err != nil {
		log.Fatal("Gemini OCR init failed:", err)
	}
	defer ocrProcessor.Close()

	submissionRepo := postgres.NewSubmissionRepo(db)
	auditRepo := postgres.NewAuditRepo(db)
	ocrService := service.NewOCRService(submissionRepo, auditRepo, minioClient, ocrProcessor)

	// Read answer.jpeg
	imageData, err := os.ReadFile("answer.jpeg")
	if err != nil {
		log.Fatal("Failed to read answer.jpeg:", err)
	}
	log.Printf("Loaded answer.jpeg (%d bytes)", len(imageData))

	// Upload to MinIO
	filename := fmt.Sprintf("test_%d.jpeg", time.Now().Unix())
	objectName := "submissions/" + filename
	url, err := minioClient.UploadFile(ctx, objectName, imageData, "image/jpeg")
	if err != nil {
		log.Fatal("Upload failed:", err)
	}
	log.Printf("✓ Uploaded to MinIO: %s (URL: %s)", objectName, url)

	// Create exam
	examID := uuid.New()
	questionID := uuid.New()
	
	_, err = db.NewInsert().Model(&map[string]interface{}{
		"id":        examID,
		"title":     "Chemistry Test",
		"subject":   "Science",
		"tenant_id": tenant.ID,
	}).Table("exams").Exec(ctx)
	if err != nil {
		log.Fatal("Exam creation failed:", err)
	}

	// Add question
	_, err = db.NewInsert().Model(&map[string]interface{}{
		"id":            questionID,
		"exam_id":       examID,
		"question_text": "Explain the Chlor-Alkali process and bleaching powder properties.",
		"points":        10,
		"answer_type":   "essay",
	}).Table("questions").Exec(ctx)
	if err != nil {
		log.Fatal("Question creation failed:", err)
	}
	log.Printf("✓ Created exam with question")

	// Create submission
	submissionID := uuid.New()
	_, err = db.NewInsert().Model(&map[string]interface{}{
		"id":                submissionID,
		"exam_id":           examID,
		"student_id":        "TEST_STUDENT_001",
		"tenant_id":         tenant.ID,
		"processing_status": "pending",
		"ocr_results":       fmt.Sprintf(`[{"page_number":1,"image_url":"%s"}]`, objectName),
	}).Table("submissions").Exec(ctx)
	if err != nil {
		log.Fatal("Submission creation failed:", err)
	}
	log.Printf("✓ Created submission: %s", submissionID)

	// Trigger OCR
	log.Println("⏳ Running OCR...")
	if err := ocrService.ProcessSubmission(ctx, submissionID); err != nil {
		log.Fatal("OCR failed:", err)
	}

	// Check result
	result, err := ocrService.GetByID(ctx, submissionID)
	if err != nil {
		log.Fatal("Failed to get submission:", err)
	}

	log.Println("\n=== OCR RESULT ===")
	if len(result.OCRResults) > 0 && result.OCRResults[0].RawText != "" {
		log.Printf("Extracted Text:\n%s\n", result.OCRResults[0].RawText)
		log.Printf("Confidence: %.2f", result.OCRResults[0].Confidence)
		log.Printf("Status: %s", result.ProcessingStatus)
	} else {
		log.Println("No text extracted")
	}

	log.Println("\n✅ Test completed successfully!")
}
