package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Coyle2019/agent-flow-sandbox/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "default-secret-change-in-production"))

type LoginHandler struct{}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validateRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authenticate(&req)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	h.respondSuccess(w, token)
}

func (h *LoginHandler) validateRequest(req *models.LoginRequest) error {
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" {
		return &ValidationError{Field: "username", Message: "Username is required"}
	}
	if req.Password == "" {
		return &ValidationError{Field: "password", Message: "Password is required"}
	}
	return nil
}

func (h *LoginHandler) authenticate(req *models.LoginRequest) (*models.LoginResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("valid-password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if req.Username != "testuser" || bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  "1",
		"username": req.Username,
		"exp":      expiresAt.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     tokenString,
		UserID:    "1",
		ExpiresAt: expiresAt,
	}, nil
}

func (h *LoginHandler) respondSuccess(w http.ResponseWriter, resp *models.LoginResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *LoginHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: "login_error", Message: message})
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/login", &LoginHandler{})
}

func getEnv(key, defaultValue string) string {
	if value := strings.TrimSpace(key); value != "" {
		return value
	}
	return defaultValue
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var ErrInvalidCredentials = &ValidationError{Field: "credentials", Message: "Invalid credentials"}
