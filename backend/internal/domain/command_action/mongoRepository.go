package command_action

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoCommandActionRepository struct {
	collection *mongo.Collection
}

func NewMongoCommandTypeRepository(col *mongo.Collection) MongoCommandActionRepository {
	return MongoCommandActionRepository{collection: col}
}

func (repo *MongoCommandActionRepository) Create(
	ctx context.Context,
	cmdAction *CommandAction,
) (*string, error) {
	cmdAction.ID = primitive.NewObjectID()
	_, err := repo.collection.InsertOne(ctx, cmdAction)
	if err != nil {
		return nil, err
	}

	idHex := cmdAction.ID.Hex()

	return &idHex, nil
}

func (repo *MongoCommandActionRepository) List(
	ctx context.Context,
	filter CommandActionFilter,
) ([]*CommandAction, error) {

	return nil, nil
}

func (repo *MongoCommandActionRepository) GetByName(
	ctx context.Context,
	name string,
) (*CommandAction, error) {
	filter := bson.D{{Key: "name", Value: name}}

	var cmdAction CommandAction
	err := repo.collection.FindOne(ctx, filter).Decode(&cmdAction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrCommandActionNotFound
		}

		return nil, ErrParsingCmd
	}

	return &cmdAction, nil
}

func (repo *MongoCommandActionRepository) Update(
	ctx context.Context,
	cmdAction *CommandAction,
) error {
	updateData, err := bson.Marshal(cmdAction)
	if err != nil {
		return err
	}

	var updateMap bson.M
	if err := bson.Unmarshal(updateData, &updateMap); err != nil {
		return err
	}

	delete(updateMap, "_id")
	delete(updateMap, "name")
	filter := bson.D{{Key: "name", Value: cmdAction.Name}}
	update := bson.D{{Key: "$set", Value: updateMap}}

	res, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrCommandActionNotFound
	}
	return nil
}

func (repo *MongoCommandActionRepository) Delete(
	ctx context.Context,
	name string,
) error {
	filter := bson.D{{Key: "name", Value: name}}

	res, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrCommandActionNotFound
	}

	return nil
}
