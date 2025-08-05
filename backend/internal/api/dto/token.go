package dto

import "time"

type NewTokenRequest struct {
}

type NewTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
