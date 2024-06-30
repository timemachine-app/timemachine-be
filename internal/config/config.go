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
	JwtSecret string
}

type ServerConfig struct {
	Port int
}

type ClientsConfig struct {
	Gemini          GeminiConfig
	OpenAI          OpenAIConfig
	SignInWithApple SignInWithAppleConfig
	Superbase       SuperbaseConfig
}

type GeminiConfig struct {
	Key   string
	Model string
}

type OpenAIConfig struct {
	Key       string
	Model     string
	MaxTokens int
}

type SignInWithAppleConfig struct {
	AppleClientId string
	TeamId        string
	KeyId         string
	PrivateKey    string
}

type SuperbaseConfig struct {
	Url              string
	Key              string
	AccountTableName string
	UsageTableName   string
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

	SearchContextHistoryPrompt           string
	SearchContextSearchTextPrompt        string
	SearchContextSystemInstructionPrompt string
	SearchContextSystemResponsePrompt    string
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
