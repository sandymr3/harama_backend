package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3-flash-preview")
	resp, err := model.GenerateContent(ctx, genai.Text("Say 'Hello, I am working!' if you can hear me."))
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Candidates) == 0 {
		log.Fatal("No candidates in response")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Println(part)
	}
}
