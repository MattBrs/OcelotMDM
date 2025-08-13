package device

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Type         string             `bson:"type"`
	Status       string             `bson:"status,omitempty"`
	IPAddress    string             `bson:"ip_address,omitempty"`
	LastSeen     int64              `bson:"last_seen,omitempty"`
	Tags         []string           `bson:"tags,omitempty"`
	Architecture string             `bson:"architecture"`
}
