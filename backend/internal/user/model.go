package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	UpdatedBy primitive.ObjectID `bson:"updated_by"`
	Enabled   bool               `bson:"enabled"`
	Admin     bool               `bson:"admin"`
}
