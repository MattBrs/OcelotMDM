package token_dto

import "time"

type NewTokenRequest struct {
}

type NewTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type NewTokenResponseErr struct {
	Error string `json:"error"`
}
