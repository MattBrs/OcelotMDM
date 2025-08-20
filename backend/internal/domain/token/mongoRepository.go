package token

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoTokenRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) MongoTokenRepository {
	return MongoTokenRepository{col}
}

func (repo MongoTokenRepository) Add(ctx context.Context, token Token) error {
	token.ID = primitive.NewObjectID()
	res, err := repo.collection.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	fmt.Println("Created new token with ID: ", res.InsertedID)
	return nil
}

func (repo MongoTokenRepository) Verify(ctx context.Context, otp string) (Token, error) {
	filter := bson.D{{Key: "token", Value: otp}}

	var token Token
	err := repo.collection.FindOneAndDelete(ctx, filter).Decode(&token)
	if err != nil {
		return token, ErrOtpNotFound
	}

	return token, nil
}
