package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/timemachine-app/timemachine-be/gemini"
	"github.com/timemachine-app/timemachine-be/internal/config"
	"github.com/timemachine-app/timemachine-be/openai"
	"github.com/timemachine-app/timemachine-be/util"
)

const (
	inputFormPhotoKey            = "timemachine-photo"
	inputFormMessageKey          = "timemachine-message"
	inputFormDate                = "timemachine-date"
	inputFormPrevTimelineEvents  = "timemachine-prev-timeline-events"
	inputFormPrevTimelineSummary = "timeline-summary"

	inputFormHistory    = "timemachine-history"
	inputFormSearchText = "timemachine-search-text"

	genericProcessingError = "Failed to process your request"
	genericBadRequestError = "Bad Input Request"
)

type EventHandler struct {
	openAIConfig config.OpenAIConfig
	geminiConfig config.GeminiConfig
	eventPrompts config.EventPromptsConfig
}

func NewEventHandler(openAIConfig config.OpenAIConfig, geminiConfig config.GeminiConfig, eventPrompts config.EventPromptsConfig) *EventHandler {
	return &EventHandler{
		openAIConfig: openAIConfig,
		geminiConfig: geminiConfig,
		eventPrompts: eventPrompts,
	}
}

func (h *EventHandler) ProcessEvent(c *gin.Context) {
	contextPrompt := ""
	timelineSummary := c.PostForm(inputFormPrevTimelineSummary)
	if timelineSummary != "" {
		contextPrompt = fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextTimelineDetailsPrompt, timelineSummary)
	}

	eventTime := c.PostForm(inputFormDate)
	if eventTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": genericBadRequestError})
		return
	}
	contextPrompt = contextPrompt + fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextTimePrompt, eventTime)

	eventMessage := c.PostForm(inputFormMessageKey)
	if eventMessage != "" {
		contextPrompt = contextPrompt +
			fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextInputMessagePrompt, eventMessage)
	}

	// previousEvents := c.PostForm(inputFormPrevTimelineEvents)
	// if previousEvents != "" {
	// 	contextPrompt = contextPrompt +
	// 		fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextPrevTimelinePrompt, previousEvents)
	// }

	// Handle file input
	var imageBytes *[]byte = nil
	file, _, err := c.Request.FormFile(inputFormPhotoKey)
	if err == nil {
		defer file.Close()
		// Read file content
		currentImageBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": genericProcessingError})
			return
		}
		imageBytes = &currentImageBytes
	}

	// response, err := openai.CallOpenAIAPI(
	// 	contextPrompt, &imageBytes,
	// 	h.eventPrompts.EventContextSystemInstructionPrompt,
	// 	h.eventPrompts.EventContextSystemResponsePrompt,
	// 	h.openAIConfig.Key,
	// 	h.openAIConfig.Model,
	// 	h.openAIConfig.MaxTokens)

	response, err := gemini.CallGeminiAPI(
		contextPrompt, imageBytes,
		h.eventPrompts.EventContextSystemInstructionPrompt,
		h.eventPrompts.EventContextSystemResponsePrompt,
		h.geminiConfig.Key,
		h.geminiConfig.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": genericProcessingError})
		return
	}

	// clean json
	cleanResponse := util.CleanLLMJson(response)

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanResponse), &jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": genericProcessingError})
		return
	}

	// Return the JSON data as a response
	c.JSON(http.StatusOK, jsonData)
}

func (h *EventHandler) Search(c *gin.Context) {
	contextPrompt := ""
	inputFormHistory := c.PostForm(inputFormHistory)
	if inputFormHistory == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": genericBadRequestError})
		return
	}
	contextPrompt = fmt.Sprintf("%s: %s. ", h.eventPrompts.SearchContextHistoryPrompt, inputFormHistory)

	searchText := c.PostForm(inputFormSearchText)
	if searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": genericBadRequestError})
		return
	}
	contextPrompt = contextPrompt + fmt.Sprintf("%s: %s. ", h.eventPrompts.SearchContextSearchTextPrompt, searchText)

	response, err := openai.CallOpenAIAPI(
		contextPrompt, nil,
		h.eventPrompts.SearchContextSystemInstructionPrompt,
		h.eventPrompts.SearchContextSystemResponsePrompt,
		h.openAIConfig.Key,
		h.openAIConfig.Model,
		h.openAIConfig.MaxTokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": genericProcessingError})
		return
	}

	// clean json
	cleanResponse := util.CleanLLMJson(response)

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanResponse), &jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": genericProcessingError})
		return
	}

	// Return the JSON data as a response
	c.JSON(http.StatusOK, jsonData)
}
