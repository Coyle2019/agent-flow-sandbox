package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Coyle2019/agent-flow-sandbox/handlers"
	"github.com/Coyle2019/agent-flow-sandbox/models"
)

func TestLoginHandler_Success(t *testing.T) {
	reqBody := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := handlers.LoginHandler{}
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp models.LoginResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Token == "" {
		t.Error("Expected non-empty token")
	}
	if resp.UserID == "" {
		t.Error("Expected non-empty user_id")
	}
	if resp.ExpiresAt.IsZero() {
		t.Error("Expected non-zero expires_at")
	}
}

func TestLoginHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler := handlers.LoginHandler{}
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestLoginHandler_MissingUsername(t *testing.T) {
	reqBody := map[string]string{"password": "password123"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := handlers.LoginHandler{}
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLoginHandler_MissingPassword(t *testing.T) {
	reqBody := map[string]string{"username": "testuser"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := handlers.LoginHandler{}
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	reqBody := models.LoginRequest{
		Username: "wronguser",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := handlers.LoginHandler{}
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var errResp models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp.Error != "login_error" {
		t.Errorf("Expected error 'login_error', got '%s'", errResp.Error)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := handlers.LoginHandler{}
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
