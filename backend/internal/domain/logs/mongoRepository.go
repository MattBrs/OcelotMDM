package logs

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoLogRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) MongoLogRepository {
	return MongoLogRepository{col}
}

func (repo MongoLogRepository) Add(ctx context.Context, log Log) error {
	res, err := repo.collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	fmt.Println("added log with ID: ", res.InsertedID)
	return nil
}
