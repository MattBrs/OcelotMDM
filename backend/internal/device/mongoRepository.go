package device

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	device.ID = primitive.NewObjectID()
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

func (repo MongoDeviceRepository) GetByName(ctx context.Context, name string) (*Device, error) {
	filter := bson.D{{Key: "name", Value: name}}
	var device Device
	err := repo.collection.FindOne(ctx, filter).Decode(&device)

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (repo MongoDeviceRepository) Update(ctx context.Context, device *Device) error {
	updateData, err := bson.Marshal(device)
	if err != nil {
		return err
	}

	var updateMap bson.M
	if err := bson.Unmarshal(updateData, &updateMap); err != nil {
		return err
	}

	delete(updateMap, "_id")
	filter := bson.D{{Key: "_id", Value: device.ID}}
	update := bson.D{{Key: "$set", Value: updateMap}}

	_, err = repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrDeviceNotUpdated
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
