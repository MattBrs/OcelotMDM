package command

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(ctx context.Context, command *Command) (*string, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*Command, error)
	Update(ctx context.Context, command *Command) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter CommandFilter) ([]*Command, error)
}
