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

// @BasePath /command_actions

// Creates a new command action
// @Summary Creates a new commmand action that is then used to create commands
// @Schemes
// @Description Creates a new commmand action that is then used to create commands. What a command does is determined by a command_action. Example: command_action='install_binary', the new command that is enqueued with this command action will install a binary.
// @Tags command_actions
// @Accept json
// @Produce json
// @Param command_action body command_action_dto.AddNewCommandActionRequest true "command action data"
// @Success 200 {object} command_action_dto.AddNewCommandActionResponse
// @Failure 400 {object} command_action_dto.ResponseErr
// @Failure 401 {object} command_action_dto.ResponseErr
// @Failure 500 {object} command_action_dto.ResponseErr
// @Router /command_actions/new [post]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
func (handler *CommandActionHandler) AddNewCommandAction(ctx *gin.Context) {
	var req command_action_dto.AddNewCommandActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			command_action_dto.ResponseErr{
				Error: err.Error(),
			},
		)
		return
	}

	cmdAct := command_action.CommandAction{
		Name:            req.Name,
		Description:     req.Description,
		RequiredOnlne:   *req.RequiredOnline,
		DefaultPriority: req.DefaultPriority,
		PayloadRequired: *req.PayloadRequired,
		TokenRequired:   *req.TokenRequired,
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

// @BasePath /command_actions

// Lists the command actions
// @Summary Returns a list of the created command_actions. With filter.
// @Schemes
// @Description  Returns a list of the created command_actions. The command actions are filtered by some attributes.
// @Tags command_actions
// @Accept json
// @Produce json
// @Param name query string false "Command action name"
// @Success 200 {object} command_action_dto.ListCommandActionResponse
// @Failure 400 {object} command_action_dto.ResponseErr
// @Failure 401 {object} command_action_dto.ResponseErr
// @Failure 500 {object} command_action_dto.ResponseErr
// @Router /command_actions/list [get]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
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

// @BasePath /command_actions

// Deletes command action
// @Summary Deletes a command action
// @Schemes
// @Description deletes a command action by name (which is unique)
// @Tags command_actions
// @Accept json
// @Produce json
// @Param command_action body command_action_dto.DeleteCommandActionRequest true "command action name"
// @Success 200 {object} command_action_dto.AddNewCommandActionResponse
// @Failure 400 {object} command_action_dto.ResponseErr
// @Failure 401 {object} command_action_dto.ResponseErr
// @Failure 500 {object} command_action_dto.ResponseErr
// @Router /command_actions/delete [post]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
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
