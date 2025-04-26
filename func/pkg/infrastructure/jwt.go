package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func ValidateJWT(token string) (string, error) {
	authURL := "http://localhost:8080/api/auth/validate"

	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return "", fmt.Errorf("could not create validation request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("authentication service unreachable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid token: status %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("could not decode auth response: %v", err)
	}

	if !authResp.Success {
		return "", fmt.Errorf("invalid token: %s", authResp.Error)
	}

	return GetUserIDFromToken(token)
}

func GetUserIDFromToken(token string) (string, error) {
	userURL := "http://localhost:8080/api/auth/user"

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return "", fmt.Errorf("could not create user request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("authentication service unreachable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user: status %d", resp.StatusCode)
	}

	var userData struct {
		Content struct {
			Email string `json:"email"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return "", fmt.Errorf("could not decode user response: %v", err)
	}

	return userData.Content.Email, nil
}
