package binary

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoBinaryRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) MongoBinaryRepository {
	return MongoBinaryRepository{
		col,
	}
}

func (repo MongoBinaryRepository) Add(ctx context.Context, binary Binary) error {
	res, err := repo.collection.InsertOne(ctx, binary)
	if err != nil {
		return err
	}

	fmt.Println("added binary with ID: ", res.InsertedID)
	return nil
}

func (repo MongoBinaryRepository) Get(
	ctx context.Context,
	binaryName string,
) (*Binary, error) {
	filter := bson.D{{Key: "name", Value: binaryName}}
	fmt.Printf("Running query: %+v\n", filter)

	var binary Binary
	err := repo.collection.FindOne(ctx, filter).Decode(&binary)

	if err != nil {
		fmt.Println("could not find binary with name: ", binaryName)
		return nil, err
	}

	return &binary, nil
}

func (repo MongoBinaryRepository) ListBinaries(
	ctx context.Context,
) ([]*Binary, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var binaries []*Binary
	for cursor.Next(ctx) {
		var b Binary
		if err = cursor.Decode(&b); err != nil {
			return nil, err
		}

		binaries = append(binaries, &b)
	}

	return binaries, nil
}
