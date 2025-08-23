package command_dto

import "github.com/MattBrs/OcelotMDM/internal/domain/command"

type AddNewCommadRequest struct {
	CommandActionName string `bson:"command_action_name"`
	DeviceName        string `bson:"device_name"`
	Payload           string `bson:"payload"`
	Priority          uint   `bson:"priority"`
}

type AddNewCommadResponse struct {
	ID                string                `bson:"id"`
	CommandActionName string                `bson:"command_action_name"`
	Status            command.CommandStatus `bson:"status"`
}

type ResponseErr struct {
	Error string `bson:"error"`
}

type ListCommandsResponse struct {
	Commands []*command.Command `bson:"commands"`
}

type DeleteCommandRequest struct {
	ID string `bson:"id"`
}

type DeleteCommandResponse struct {
	ID string `bson:"id"`
}
