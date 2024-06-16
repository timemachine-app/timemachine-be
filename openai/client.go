package openai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	endpoint = "https://api.openai.com/v1/chat/completions"
)

// Content represents content in a message
type Content struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents the URL of an image
type ImageURL struct {
	URL string `json:"url"`
}

// Message represents a message in a conversation
type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

// Payload represents the payload sent to the OpenAI API
type Payload struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// OpenAIResponse represents the structure of the response from the OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// CallOpenAIAPI calls the OpenAI API for image processing
func CallOpenAIAPI(
	contextPrompt string,
	imageBytes *[]byte,
	systemPrompt string,
	responsePrompt string,
	apiKey string,
	model string,
	maxTokens int) (string, error) {
	var messages []Message

	// systemPrompt
	if systemPrompt != "" {
		messages = append(messages, Message{
			Role: "system",
			Content: []Content{
				{
					Type: "text",
					Text: systemPrompt,
				},
			},
		})
	}

	// contextPrompt with imageBytes
	var promptContents []Content
	if imageBytes != nil {
		encodedImage := base64.StdEncoding.EncodeToString(*imageBytes)
		promptContents = append(promptContents, Content{
			Type: "image_url",
			ImageURL: ImageURL{
				URL: "data:image/jpeg;base64," + encodedImage,
			},
		})
	}
	// Add the text content
	promptContents = append(promptContents, Content{
		Type: "text",
		Text: contextPrompt,
	})
	messages = append(messages, Message{
		Role:    "user",
		Content: promptContents,
	})

	// responsePrompt
	if responsePrompt != "" {
		messages = append(messages, Message{
			Role: "system",
			Content: []Content{
				{
					Type: "text",
					Text: responsePrompt,
				},
			},
		})
	}

	payload := Payload{
		Model:     model,
		Messages:  messages,
		MaxTokens: maxTokens,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %w", err)
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var openAIResponse OpenAIResponse
	if err := json.Unmarshal(responseData, &openAIResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	if len(openAIResponse.Choices) == 0 {
		// return "", fmt.Errorf("no choices in the response")
		return string(responseData), nil
	}

	return openAIResponse.Choices[0].Message.Content, nil
}
