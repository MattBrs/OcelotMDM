package command_action_dto

import "github.com/MattBrs/OcelotMDM/internal/domain/command_action"

type AddNewCommandActionRequest struct {
	Name            string `bson:"name" binding:"required"`
	Description     string `bson:"description" binding:"required"`
	RequiredOnline  bool   `bson:"required_online" binding:"required"`
	DefaultPriority uint   `bson:"default_priority" binding:"required"`
}

type AddNewCommandActionResponse struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
}

type DeleteCommandActionRequest struct {
	Name string `bson:"name" binding:"required"`
}

type DeleteCommandActionResponse struct {
	Status string `bson:"status"`
}

type ListCommandActionResponse struct {
	CommandActions []*command_action.CommandAction `bson:"command_actions"`
}

type ResponseErr struct {
	Error string `bson:"error"`
}
