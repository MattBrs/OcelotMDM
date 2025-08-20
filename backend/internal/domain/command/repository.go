package command

import "context"

type Repository interface {
	Create(ctx context.Context, command *Command) error
	GetById(ctx context.Context, id string) (*Command, error)
	Update(ctx context.Context, command *Command) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter CommandFilter) ([]*Command, error)
}
