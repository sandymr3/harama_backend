package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"harama/internal/ocr"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set in .env")
	}

	// Read the image file
	imgBytes, err := os.ReadFile("../../answer.jpeg")
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}

	// Initialize OCR processor
	processor, err := ocr.NewGeminiOCRProcessor(apiKey)
	if err != nil {
		log.Fatalf("Failed to create processor: %v", err)
	}
	defer processor.Close()

	fmt.Println("🚀 Processing OCR for answer.jpeg...")

	// Extract text
	result, err := processor.ExtractText(context.Background(), imgBytes, "image/jpeg")
	if err != nil {
		log.Fatalf("OCR failed: %v", err)
	}

	fmt.Println("\n--- OCR OUTPUT ---")
	fmt.Println(result.RawText)
	fmt.Println("------------------")
	fmt.Printf("\nConfidence: %.2f\n", result.Confidence)
}
