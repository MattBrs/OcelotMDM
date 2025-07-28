package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoConnection struct {
	client *mongo.Client
}

func NewMongoConnection(config DbConfig) (MongoConnection, error) {
	if err := config.verifyMongoValidity(); err != nil {
		return MongoConnection{nil}, err
	}

	mongoConnectionStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=%s",
		config.Username,
		config.Password,
		config.ClusterURL,
		config.AppName,
	)

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoConnectionStr).SetServerAPIOptions(serverApi)

	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		return MongoConnection{nil}, err
	}

	return MongoConnection{mongoClient}, nil
}

func (conn MongoConnection) CloseMongoConnection() error {
	return conn.client.Disconnect(context.TODO())
}

func (conn MongoConnection) Ping() error {
	return conn.client.Ping(context.TODO(), readpref.Primary())
}

func (conn MongoConnection) GetCollection(dbName string, collectionName string) *mongo.Collection {
	return conn.client.Database(dbName).Collection(collectionName)
}
