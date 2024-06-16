package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

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

	genericProcessingError = "Failed to process an event"
	genericBadRequestError = "Bad Input Request"
)

type EventHandler struct {
	openAIConfig config.OpenAIConfig
	eventPrompts config.EventPromptsConfig
}

func NewEventHandler(openAIConfig config.OpenAIConfig, eventPrompts config.EventPromptsConfig) *EventHandler {
	return &EventHandler{
		openAIConfig: openAIConfig,
		eventPrompts: eventPrompts,
	}
}

func (h *EventHandler) ProcessEvent(c *gin.Context) {
	contextPrompt := ""
	timelineSummary := c.PostForm(inputFormPrevTimelineSummary)
	if timelineSummary == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": genericBadRequestError})
		return
	}
	contextPrompt = fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextTimelineDetailsPrompt, timelineSummary)

	eventTime := c.PostForm(inputFormDate)
	if eventTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": genericBadRequestError})
		return
	}
	contextPrompt = fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextTimePrompt, eventTime)

	eventMessage := c.PostForm(inputFormMessageKey)
	if eventMessage != "" {
		contextPrompt = contextPrompt +
			fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextInputMessagePrompt, eventMessage)
	}

	previousEvents := c.PostForm(inputFormPrevTimelineEvents)
	if previousEvents != "" {
		contextPrompt = contextPrompt +
			fmt.Sprintf("%s: %s. ", h.eventPrompts.EventContextPrevTimelinePrompt, previousEvents)
	}

	file, _, err := c.Request.FormFile(inputFormPhotoKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": genericBadRequestError})
		return
	}
	defer file.Close()

	// Read file content
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": genericBadRequestError})
		return
	}

	response, err := openai.CallOpenAIAPI(
		contextPrompt, &imageBytes,
		h.eventPrompts.EventContextSystemInstructionPrompt,
		h.eventPrompts.EventContextSystemResponsePrompt,
		h.openAIConfig.Key,
		h.openAIConfig.Model,
		h.openAIConfig.MaxTokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": genericProcessingError})
		return
	}

	// clean json
	response = util.CleanOpenAIJson(response)

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": genericProcessingError})
		return
	}

	// Return the JSON data as a response
	c.JSON(http.StatusOK, jsonData)
}
