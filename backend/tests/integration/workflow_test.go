package integration

// This test verifies the end-to-end workflow of the backend:
// 1. OCR processing of an answer image (using Gemini Vision)
// 2. Grading of the extracted text (using Gemini 1.5 Pro/Flash)
// It ensures that the core grading engine logic, including multi-evaluator consensus and partial credit calculation, functions correctly.

import (
	"context"
	"os"
	"testing"
	"time"

	"harama/internal/ai/gemini"
	"harama/internal/domain"
	"harama/internal/grading"
	"harama/internal/ocr"

	"github.com/google/uuid"
)

func TestOCRAndGradingWorkflow(t *testing.T) {
	// 1. Setup Environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 2. Initialize Services
	t.Log("Initializing services...")
	
	// AI Client
	aiClient, err := gemini.NewClient(apiKey)
	if err != nil {
		t.Fatalf("Failed to create Gemini client: %v", err)
	}

	// OCR Processor
	ocrProcessor, err := ocr.NewGeminiOCRProcessor(apiKey)
	if err != nil {
		t.Fatalf("Failed to create OCR processor: %v", err)
	}
	defer ocrProcessor.Close()

	// Grading Engine
	gradingEngine := grading.NewEngine(aiClient)

	// 3. Load Image
	t.Log("Reading answer.jpeg...")
	imgBytes, err := os.ReadFile("../../answer.jpeg")
	if err != nil {
		t.Fatalf("Failed to read answer.jpeg: %v", err)
	}

	// 4. Perform OCR
	t.Log("Running OCR...")
	ocrResult, err := ocrProcessor.ExtractText(ctx, imgBytes, "image/jpeg")
	if err != nil {
		t.Fatalf("OCR failed: %v", err)
	}
	t.Logf("OCR Text Extracted: \n%s", ocrResult.RawText)

	// 5. Setup Grading Context (Rubric & Question)
	questionText := "Explain the Chlor-Alkali process and the preparation of Bleaching Powder. Include chemical equations."
	subject := "Chemistry"

	rubric := domain.Rubric{
		ID: uuid.New(),
		QuestionID: uuid.New(),
		FullCreditCriteria: []domain.Criterion{
			{ID: "c1", Description: "Mentions Chlor-Alkali process reactants (NaCl and H2O) and products (NaOH, Cl2, H2)", Points: 2.0},
			{ID: "c2", Description: "Correct chemical equation for Chlor-Alkali process: 2NaCl + 2H2O -> 2NaOH + Cl2 + H2", Points: 2.0},
			{ID: "c3", Description: "Mentions Bleaching Powder formula is CaOCl2", Points: 1.0},
			{ID: "c4", Description: "Explains preparation from Slaked Lime (Ca(OH)2) and Chlorine", Points: 2.0},
			{ID: "c5", Description: "Correct chemical equation for Bleaching Powder: Ca(OH)2 + Cl2 -> CaOCl2 + H2O", Points: 3.0},
		},
		StrictMode: false,
	}

	// Construct Answer Segment from OCR result
	answer := domain.AnswerSegment{
		ID: uuid.New(),
		Text: ocrResult.RawText,
	}

	// 6. Perform Grading
	t.Log("Running Grading Engine...")
	finalGrade, multiEval, err := gradingEngine.GradeAnswer(ctx, answer, rubric, subject, questionText)
	if err != nil {
		t.Fatalf("Grading failed: %v", err)
	}

	// 7. Validate Results
	t.Logf("--- GRADING RESULTS ---")
	t.Logf("Final Score: %.2f / 10.0", finalGrade.FinalScore)
	t.Logf("Confidence: %.2f", finalGrade.Confidence)
	t.Logf("Status: %s", finalGrade.Status)
	t.Logf("Reasoning: %s", finalGrade.Reasoning)
	
	if multiEval != nil {
		t.Logf("--- MULTI-EVALUATOR DETAILS ---")
		for _, eval := range multiEval.Evaluations {
			t.Logf("[%s] Score: %.2f (Conf: %.2f)", eval.AIEvaluatorID, eval.Score, eval.Confidence)
			t.Logf("  Criteria Met: %+v", eval.CriteriaMet)
		}
		t.Logf("Variance: %.4f", multiEval.Variance)
	}

	if finalGrade.FinalScore < 5.0 {
		t.Error("Expected a passing score (> 5.0) given the answer quality visible in the image.")
	}
	if finalGrade.Confidence < 0.5 {
		t.Log("Confidence is low. Check if OCR text was garbled or rubric is too strict.")
	}
}