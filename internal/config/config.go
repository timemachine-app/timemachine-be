package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Clients   ClientsConfig
	Prompts   PromptsConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port int
}

type ClientsConfig struct {
	Gemini GeminiConfig
	OpenAI OpenAIConfig
}

type GeminiConfig struct {
	Key string
}

type OpenAIConfig struct {
	Key       string
	Model     string
	MaxTokens int
}

type PromptsConfig struct {
	EventPrompts EventPromptsConfig
}

type EventPromptsConfig struct {
	EventContextTimelineDetailsPrompt   string
	EventContextTimePrompt              string
	EventContextInputMessagePrompt      string
	EventContextPrevTimelinePrompt      string
	EventContextSystemInstructionPrompt string
	EventContextSystemResponsePrompt    string
}

type RateLimitConfig struct {
	RateLimit   int
	WindowInSec int64
}

func LoadConfig(configName string) (*Config, error) {
	var config Config

	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
