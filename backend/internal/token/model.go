package token

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Token       string             `bson:"token"`
	ExpiresAt   time.Time          `bson:"expires_at"`
	RequestedBy string             `bson:"requested_by"`
}
