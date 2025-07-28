package device

import "go.mongodb.org/mongo-driver/v2/mongo"

// TODO: add functions defined in internal/device/repository.go to expose
// all the proper functions

type MongoDeviceRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) MongoDeviceRepository {
	return MongoDeviceRepository{col}
}


