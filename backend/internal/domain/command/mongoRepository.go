package command

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoCommandRepository struct {
	collection *mongo.Collection
}

// TODO: finish implementation of dao

func NewMongoRepository(col *mongo.Collection) MongoCommandRepository {
	return MongoCommandRepository{col}
}

func Create(ctx context.Context, command *Command) error {
	return nil
}

func GetById(ctx context.Context, id string) (*Command, error) {
	return nil, nil
}

func Update(ctx context.Context, command *Command) error {
	return nil
}

func Delete(ctx context.Context, id string) error {
	return nil
}

func List(ctx context.Context, filter CommandFilter) ([]*Command, error) {
	return nil, nil
}
