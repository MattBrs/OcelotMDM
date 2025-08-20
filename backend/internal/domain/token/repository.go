package token

import "context"

type Repository interface {
	Add(ctx context.Context, token Token) error
	Verify(ctx context.Context, token string) (Token, error)
}
