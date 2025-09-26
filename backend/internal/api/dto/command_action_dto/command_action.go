package command_action_dto

import "github.com/MattBrs/OcelotMDM/internal/domain/command_action"

type AddNewCommandActionRequest struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	RequiredOnline  *bool  `json:"required_online" binding:"required"`
	DefaultPriority uint   `json:"default_priority" binding:"required"`
	PayloadRequired *bool  `json:"payload_required" binding:"required"`
	TokenRequired   *bool  `json:"token_required" binding:"required"`
}

type AddNewCommandActionResponse struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name"`
}

type DeleteCommandActionRequest struct {
	Name string `json:"name" binding:"required"`
}

type DeleteCommandActionResponse struct {
	Status string `json:"status"`
}

type ListCommandActionResponse struct {
	CommandActions []*command_action.CommandAction `json:"command_actions"`
}

type ResponseErr struct {
	Error string `json:"error"`
}
