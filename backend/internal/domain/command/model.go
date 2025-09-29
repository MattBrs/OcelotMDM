package command

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommandStatus struct {
	Id     int    `bson:"status_id"`
	Status string `bson:"status_value"`
}

func StatusFromString(val string) *CommandStatus {
	switch val {
	case WAITING.Status:
		return &COMPLETED
	case QUEUED.Status:
		return &QUEUED
	case ACKED.Status:
		return &ACKED
	case COMPLETED.Status:
		return &COMPLETED
	case ERRORED.Status:
		return &ERRORED
	default:
		return nil
	}
}

var (
	WAITING   = CommandStatus{Id: 1, Status: "waiting"}
	QUEUED    = CommandStatus{Id: 2, Status: "queued"}
	ACKED     = CommandStatus{Id: 3, Status: "acknowledged"}
	COMPLETED = CommandStatus{Id: 4, Status: "completed"}
	ERRORED   = CommandStatus{Id: 5, Status: "errored"}
)

type Command struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	CommandActionName string             `bson:"command_action_name"`
	DeviceName        string             `bson:"device_name"`
	Payload           string             `bson:"payload,omitempty"`
	Status            CommandStatus      `bson:"status"`
	CreatedAt         *time.Time         `bson:"created_at"`
	QueuedAt          *time.Time         `bson:"queued_at"`
	CompletedAt       *time.Time         `bson:"completed_at"`
	Priority          uint               `bson:"priority"`
	RequestedBy       string             `bson:"requested_by"`
	ErrorDescription  string             `bson:"error_desc,omitempty"`
	QueueID           primitive.ObjectID `bson:"queue_id,omitempty"`
	RequiredOnline    bool               `bson:"required_online,omitempty"`
	TokenRequired     bool               `bson:"token_required,omitempty"`
	Data              string             `bson:"data,omitempty"`
	CallbackURL       *string            `bson:"callback_url,omitempty"`
	CallbackSecret    *string            `bson:"callback_secret,omitempty"`
}
