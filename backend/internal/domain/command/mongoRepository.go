package command

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoCommandRepository struct {
	collection *mongo.Collection
}

// TODO: finish implementation of dao

func NewMongoRepository(col *mongo.Collection) MongoCommandRepository {
	return MongoCommandRepository{col}
}

func (r *MongoCommandRepository) Create(ctx context.Context, command *Command) (*string, error) {
	command.Id = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, command)
	if err != nil {
		return nil, err
	}

	idStr := command.Id.Hex()
	return &idStr, nil
}

func (r *MongoCommandRepository) GetById(ctx context.Context, id primitive.ObjectID) (*Command, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	var command Command
	err := r.collection.FindOne(ctx, filter).Decode(&command)
	if err != nil {
		var returnErr error
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			returnErr = ErrDeviceNotFound
		default:
			returnErr = ErrParsingResult
		}

		return nil, returnErr
	}

	return &command, nil
}

func (r *MongoCommandRepository) Update(ctx context.Context, command *Command) error {
	updateData, err := bson.Marshal(command)
	if err != nil {
		return err
	}

	var updateMap bson.M
	if err := bson.Unmarshal(updateData, &updateMap); err != nil {
		return err
	}

	delete(updateMap, "_id")
	filter := bson.D{{Key: "_id", Value: command.Id}}
	update := bson.D{{Key: "$set", Value: updateMap}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrCommandNotFound
	}

	return nil
}

func (r *MongoCommandRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}

	res := r.collection.FindOneAndDelete(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return ErrCommandNotFound
		}

		return res.Err()
	}

	return nil
}

func (r *MongoCommandRepository) List(ctx context.Context, filter CommandFilter) ([]*Command, error) {
	return nil, nil
}
