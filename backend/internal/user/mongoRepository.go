package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) MongoUserRepository {
	return MongoUserRepository{col}
}

func (repo *MongoUserRepository) Create(ctx context.Context, user *User) (*string, error) {
	user.ID = primitive.NewObjectID()
	res, err := repo.collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("insertOne error:", err.Error())
		return nil, err
	}

	oid := fmt.Sprintf("%s", res.InsertedID)

	return &oid, nil
}

func (repo *MongoUserRepository) List(ctx context.Context, filter UserFilter) ([]*User, error) {
	mongoFilter := bson.M{}
	if filter.Id != "" {
		mongoFilter["_id"] = filter.Id
	}

	if filter.Username != "" {
		mongoFilter["username"] = filter.Username
	}

	if filter.Admin != nil {
		mongoFilter["admin"] = filter.Admin
	}

	if filter.Enabled != nil {
		mongoFilter["enabled"] = filter.Enabled
	}

	cursor, err := repo.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}

	var users []*User
	for cursor.Next(ctx) {
		var d User
		if err = cursor.Decode(&d); err != nil {
			return nil, err
		}

		users = append(users, &d)
	}

	return users, nil
}

func (repo *MongoUserRepository) Update(ctx context.Context, user *User) error {
	updateData, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	var updateMap bson.M
	if err := bson.Unmarshal(updateData, &updateMap); err != nil {
		return err
	}

	delete(updateMap, "_id")
	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: updateMap}}

	_, err = repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrUserNotUpdated
	}

	return nil
}

func (repo *MongoUserRepository) GetById(ctx context.Context, id string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	var user User
	err = repo.collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (repo *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	filter := bson.D{{Key: "username", Value: username}}

	var user User
	err := repo.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}
