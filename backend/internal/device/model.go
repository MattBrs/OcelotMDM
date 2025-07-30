package device

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Type      string             `bson:"type"`
	Status    string             `bson:"status"`
	IPAddress string             `bson:"ip_address"`
	LastSeen  int64              `bson:"last_seen"`
	Tags      []string           `bson:"tags,omitempty"`
}
