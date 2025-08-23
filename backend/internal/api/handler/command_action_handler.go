package api

import (
	"github.com/MattBrs/OcelotMDM/internal/domain/command_action"
	"github.com/gin-gonic/gin"
)

type CommandActionHandler struct {
	service *command_action.Service
}

func NewCommandActionHandler(
	service *command_action.Service,
) *CommandActionHandler {
	return &CommandActionHandler{
		service: service,
	}
}

func (handler *CommandActionHandler) AddNewCommandAction(ctx *gin.Context) {

}

func (handler *CommandActionHandler) ListCommandActions(ctx *gin.Context) {

}

func (handler *CommandActionHandler) UpdateCommandAction(ctx *gin.Context) {

}

func (handler *CommandActionHandler) DeleteCommandAction(ctx *gin.Context) {

}
