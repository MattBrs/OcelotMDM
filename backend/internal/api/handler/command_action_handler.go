package api

import (
	"errors"
	"net/http"

	"github.com/MattBrs/OcelotMDM/internal/api/dto/command_action_dto"
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
	var req command_action_dto.AddNewCommandActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			command_action_dto.ResponseErr{
				Error: "could not parse json",
			},
		)
		return
	}

	cmdAct := command_action.CommandAction{
		Name:            req.Name,
		Description:     req.Description,
		RequiredOnlne:   req.RequiredOnline,
		DefaultPriority: req.DefaultPriority,
		PayloadRequired: req.PayloadRequired,
	}

	id, err := handler.service.AddCommandAction(ctx, &cmdAct)
	if err != nil {
		errRes := command_action_dto.ResponseErr{
			Error: "generic error",
		}
		httpStatus := http.StatusInternalServerError

		switch {
		case errors.Is(err, command_action.ErrCommandActionNameTaken):
			errRes.Error = err.Error()
			httpStatus = http.StatusConflict
		case errors.Is(err, command_action.ErrNameEmpty):
			errRes.Error = err.Error()
			httpStatus = http.StatusBadRequest
		case errors.Is(err, command_action.ErrDescriptionEmpty):
			errRes.Error = err.Error()
			httpStatus = http.StatusBadRequest
		}

		ctx.JSON(httpStatus, errRes)
		return
	}

	ctx.JSON(
		http.StatusCreated, command_action_dto.AddNewCommandActionResponse{
			ID:   *id,
			Name: cmdAct.Name,
		},
	)
}

func (handler *CommandActionHandler) ListCommandActions(ctx *gin.Context) {
	name := ctx.Query("name")

	filter := command_action.CommandActionFilter{
		Name: name,
	}

	cmdActs, err := handler.service.List(ctx, filter)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			command_action_dto.ResponseErr{
				Error: "could not parse list",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		command_action_dto.ListCommandActionResponse{
			CommandActions: cmdActs,
		},
	)
}

func (handler *CommandActionHandler) DeleteCommandAction(ctx *gin.Context) {
	var req command_action_dto.DeleteCommandActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			command_action_dto.ResponseErr{
				Error: "could not parse json",
			},
		)
		return
	}

	err := handler.service.Delete(ctx, req.Name)
	if err != nil {
		errRes := command_action_dto.ResponseErr{
			Error: "generic error",
		}
		httpStatus := http.StatusInternalServerError

		switch {
		case errors.Is(err, command_action.ErrCommandActionNotFound):
			errRes.Error = err.Error()
			httpStatus = http.StatusNotFound
		case errors.Is(err, command_action.ErrNameEmpty):
			errRes.Error = err.Error()
			httpStatus = http.StatusBadRequest
		}

		ctx.JSON(httpStatus, errRes)
		return
	}

	ctx.JSON(
		http.StatusOK,
		command_action_dto.DeleteCommandActionResponse{
			Status: "deleted",
		},
	)
}
