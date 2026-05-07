package models

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
