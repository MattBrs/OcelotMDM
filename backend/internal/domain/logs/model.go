package logs

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	deviceName       string             `bson:"device_name,omitempty"`
	RegistrationTime time.Time          `bson:"registration_time,omitempty"`
	LogData          []byte             `bson:"log_data,omitempty"`
	LogSize          int                `bson:"log_size,omitempty"`
}
