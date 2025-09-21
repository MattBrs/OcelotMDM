package binary

import "go.mongodb.org/mongo-driver/bson/primitive"

type Binary struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Data         []byte             `bson:"data,omitempty"`
	Architecture string             `bson:"architecture,omitempty"`
	Version      string             `bson:"version,omitempty"`
}
