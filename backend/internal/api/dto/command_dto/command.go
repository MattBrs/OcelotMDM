package command_dto

import "github.com/MattBrs/OcelotMDM/internal/domain/command"

type AddNewCommadRequest struct {
	CommandActionName string `json:"command_action_name" binding:"required"`
	DeviceName        string `json:"device_name" binding:"required"`
	Payload           string `json:"payload"`
	Priority          uint   `json:"priority" binding:"required"`
}

type AddNewCommadResponse struct {
	ID                string                `json:"id"`
	CommandActionName string                `json:"command_action_name"`
	Status            command.CommandStatus `json:"status"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

type ListCommandsResponse struct {
	Commands []*command.Command `json:"commands"`
}

type DeleteCommandRequest struct {
	ID string `json:"id" binding:"required"`
}

type DeleteCommandResponse struct {
	ID string `json:"id"`
}

type UpdateCommandStatusRequest struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type UpdateCommandStatusResponse struct {
	ID        string `json:"id"`
	NewStatus string `json:"new_status"`
}
