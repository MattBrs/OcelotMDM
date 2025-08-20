package token

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GenerateNewToken(ctx context.Context) (*Token, error) {
	newToken, err := s.generateToken(8, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	err = s.repo.Add(ctx, newToken)
	if err != nil {
		return nil, err
	}

	return &newToken, nil
}

func (s *Service) Verify(ctx context.Context, otp string) (bool, error) {
	token, err := s.repo.Verify(ctx, otp)
	if err != nil {
		return false, err
	}

	return token.ExpiresAt.After(time.Now()), nil
}

func (s *Service) generateToken(length int, ttl time.Duration) (Token, error) {
	buf := make([]byte, length)

	_, err := rand.Read(buf)
	if err != nil {
		return Token{}, err
	}

	rawToken := base32.
		StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(buf)

	generatedToken := Token{
		Token:       rawToken,
		ExpiresAt:   time.Now().Add(ttl),
		RequestedBy: "",
	}

	return generatedToken, nil
}
