package command

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommandStatus struct {
	Id     int    `bson:"status_id"`
	Status string `bson:"status_value"`
}

var (
	WAITING   = CommandStatus{Id: 1, Status: "waiting"}
	QUEUED    = CommandStatus{Id: 1, Status: "queued"}
	COMPLETED = CommandStatus{Id: 1, Status: "completed"}
	ERRORED   = CommandStatus{Id: 1, Status: "errored"}
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
