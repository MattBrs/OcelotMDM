package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*string, error)
	List(ctx context.Context, filter UserFilter) ([]*User, error)
	Update(ctx context.Context, user *User) error
	GetById(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
