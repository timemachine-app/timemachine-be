package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// CallGeminiAPI calls the Gemini API for image processing
func CallGeminiAPI(
	contextPrompt string,
	imageBytes *[]byte,
	systemPrompt string,
	responsePrompt string,
	apiKey string,
	model string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	genModel := client.GenerativeModel(model)

	prompt := []genai.Part{
		genai.Text(systemPrompt + "\n" + contextPrompt + "\n" + responsePrompt),
	}
	if imageBytes != nil {
		prompt = []genai.Part{
			genai.ImageData("jpeg", *imageBytes),
			genai.Text(systemPrompt + "\n" + contextPrompt + "\n" + responsePrompt),
		}
	}
	resp, err := genModel.GenerateContent(ctx, prompt...)

	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("gemini error")
	}

	println("here")

	geminiResponseText := ""

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				geminiResponseText = geminiResponseText + fmt.Sprintf("%s", part)
			}
		}
	}

	return geminiResponseText, nil
}
