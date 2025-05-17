package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type AuthResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type UserResponse struct {
	Content struct {
		ID string `json:"id"`
	} `json:"content"`
}

func ValidateJWT(tokenString string) (string, error) {
	if tokenString == "" {
		log.Println("[JWTValidation] ERROR: Empty token provided")
		return "", fmt.Errorf("empty token")
	}

	log.Printf("[JWTValidation] Starting validation for token (first 10 chars): %.10s...", tokenString)

	userID, serviceErr := validateWithAuthService(tokenString)
	if serviceErr == nil {
		return userID, nil
	}

	log.Printf("[JWTValidation] Auth service validation failed, falling back to JWT parsing: %v", serviceErr)

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Printf("[JWTValidation] ERROR parsing token: %v", err)
		return "", fmt.Errorf("invalid token format")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["userId"].(string); ok && userID != "" {
			log.Printf("[JWTValidation] Got userID from JWT: %s", userID)
			return userID, nil
		}

		if userID, ok := claims["userId"].(float64); ok {
			log.Printf("[JWTValidation] Got numeric userID from JWT: %.0f", userID)
			return fmt.Sprintf("%.0f", userID), nil
		}

		if sub, ok := claims["sub"].(string); ok {
			log.Printf("[JWTValidation] Using 'sub' claim as fallback: %s", sub)
			return sub, nil
		}
	}

	log.Println("[JWTValidation] ERROR: No valid user identifier found in token")
	return "", fmt.Errorf("no valid user identifier in token")
}

func validateWithAuthService(token string) (string, error) {
	authURL := "http://localhost:8080/api/auth/validate"
	log.Printf("[JWTValidation] Making request to auth service: %s", authURL)

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

	userURL := "http://localhost:8080/api/auth/user"
	req, err = http.NewRequest("GET", userURL, nil)
	if err != nil {
		return "", fmt.Errorf("could not create user request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("authentication service unreachable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user: status %d", resp.StatusCode)
	}

	var userData UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return "", fmt.Errorf("could not decode user response: %v", err)
	}

	if userData.Content.ID != "" {
		return userData.Content.ID, nil
	}

	return "", fmt.Errorf("empty user ID in response")
}
