package command

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommandStatus string

const (
	WAITING   CommandStatus = "waiting"
	QUEUED    CommandStatus = "queued"
	COMPLETED CommandStatus = "completed"
	ERRORED   CommandStatus = "errored"
)

type Command struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	CommandTypeId    primitive.ObjectID `bson:"command_type_id"`
	DeviceName       string             `bson:"device_name"`
	Payload          string             `bson:"payload,omitempty"`
	Status           CommandStatus      `bson:"status"`
	CreatedAt        time.Time          `bson:"created_at"`
	QueuedAt         time.Time          `bson:"queued_at"`
	CompletedAt      time.Time          `bson:"completed_at"`
	Priority         uint               `bson:"priority"`
	RequestedBy      string             `bson:"requested_by"`
	ErrorDescription string             `bson:"error_desc,omitempty"`
}
