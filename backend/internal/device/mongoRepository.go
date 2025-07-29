package device

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDeviceRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) MongoDeviceRepository {
	return MongoDeviceRepository{col}
}

func (repo MongoDeviceRepository) Create(ctx context.Context, device *Device) error {
	// doc, err := bson.Marshal(device)
	// if err != nil {
	// 	return err
	// }

	res, err := repo.collection.InsertOne(ctx, device)
	if err != nil {
		return err
	}

	fmt.Println("new device inserted with id: ", res.InsertedID)
	return nil
}

func (repo MongoDeviceRepository) GetByID(ctx context.Context, id string) (*Device, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var device Device
	err := repo.collection.FindOne(ctx, filter).Decode(&device)

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (repo MongoDeviceRepository) Update(ctx context.Context, device *Device) error {
	filter := bson.D{{Key: "_id", Value: device.ID}}
	res := repo.collection.FindOneAndUpdate(ctx, filter, device)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (repo MongoDeviceRepository) Delete(ctx context.Context, id string) error {
	filter := bson.D{{Key: "_id", Value: id}}

	res := repo.collection.FindOneAndDelete(ctx, filter)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (repo MongoDeviceRepository) List(ctx context.Context, filter DeviceFilter) ([]*Device, error) {
	mongoFilter := bson.M{}
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}

	if filter.Id != "" {
		mongoFilter["id"] = filter.Id
	}

	if filter.Name != "" {
		mongoFilter["name"] = filter.Name
	}

	if filter.Limit > 0 {
		mongoFilter["limit"] = filter.Limit
	}

	cursor, err := repo.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	for cursor.Next(ctx) {
		var d Device
		if err = cursor.Decode(&d); err != nil {
			return nil, err
		}

		devices = append(devices, &d)
	}

	return devices, nil
}
