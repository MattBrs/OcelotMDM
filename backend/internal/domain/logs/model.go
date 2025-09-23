package logs

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	DeviceName       string             `bson:"device_name"`
	RegistrationTime time.Time          `bson:"registration_time"`
	LogData          string             `bson:"log_data"`
	LogSize          int                `bson:"log_size"`
}
