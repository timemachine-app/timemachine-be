package superbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/timemachine-app/timemachine-be/internal/config"
)

type User struct {
	UserId         string `json:"UserId,omitempty"` // omit empty to exclude from POST requests
	Email          string `json:"Email"`
	ExternalUserId string `json:"ExternalUserId"`
	CreatedAt      string `json:"created_at,omitempty"` // omit empty to exclude from POST requests
}

type UsageEvent struct {
	UserId    string `json:"UserId"`
	EventType string `json:"EventType"`
}

type SupabaseClient struct {
	superbaseConfig config.SuperbaseConfig
}

func NewSupabaseClient(superbaseConfig config.SuperbaseConfig) *SupabaseClient {
	return &SupabaseClient{
		superbaseConfig: superbaseConfig,
	}
}

func (s *SupabaseClient) AddUser(user User) error {
	url := fmt.Sprintf("%s/rest/v1/%s", s.superbaseConfig.Url, s.superbaseConfig.AccountTableName)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.superbaseConfig.Key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.superbaseConfig.Key))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add user, status code: %d, response")
	}

	return nil
}

func (s *SupabaseClient) GetUser(externalUserId string) (User, error) {
	var users []User
	url := fmt.Sprintf("%s/rest/v1/%s?ExternalUserId=eq.%s", s.superbaseConfig.Url, s.superbaseConfig.AccountTableName, externalUserId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return User{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", s.superbaseConfig.Key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.superbaseConfig.Key))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return User{}, fmt.Errorf("failed to get user, status code: %d, response: %s", resp.StatusCode, bodyString)
	}

	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return User{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(users) == 0 {
		return User{}, fmt.Errorf("user not found")
	}

	return users[0], nil
}

func (s *SupabaseClient) DeleteUser(userId string) error {
	url := fmt.Sprintf("%s/rest/v1/%s?UserId=eq.%s", s.superbaseConfig.Url, s.superbaseConfig.AccountTableName, userId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", s.superbaseConfig.Key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.superbaseConfig.Key))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("failed to delete user, status code: %d, response: %s", resp.StatusCode, bodyString)
	}

	return nil
}

func (s *SupabaseClient) AddUsageEvent(event UsageEvent) error {
	url := fmt.Sprintf("%s/rest/v1/%s", s.superbaseConfig.Url, s.superbaseConfig.UsageTableName)

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal usage event data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.superbaseConfig.Key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.superbaseConfig.Key))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("failed to add usage event, status code: %d, response: %s", resp.StatusCode, bodyString)
	}

	return nil
}
