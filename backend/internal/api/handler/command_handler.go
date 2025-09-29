package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/dto/command_dto"
	"github.com/MattBrs/OcelotMDM/internal/domain/command"
	"github.com/MattBrs/OcelotMDM/internal/domain/command_action"
	"github.com/MattBrs/OcelotMDM/internal/domain/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommandHandler struct {
	service *command.Service
}

func NewCommandHandler(service *command.Service) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}

// @BasePath /command

// Enqueues a new command to target device
// @Summary Creates a new commmand that is enqueued to the device
// @Schemes
// @Description Lets a user create a command that will be enqueued as soon as possible to the mqtt topic of the target device
// @Tags command
// @Accept json
// @Produce json
// @Param command body command_dto.AddNewCommadRequest true "command data"
// @Success 200 {object} command_dto.AddNewCommadResponse
// @Failure 400 {object} command_dto.ResponseErr
// @Failure 401 {object} command_dto.ResponseErr
// @Failure 500 {object} command_dto.ResponseErr
// @Router /command/new [post]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
func (handler *CommandHandler) AddNewCommand(ctx *gin.Context) {
	var req command_dto.AddNewCommadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("could not parse json: ", err.Error())
		ctx.JSON(
			http.StatusBadRequest,
			command_dto.ResponseErr{
				Error: "could not parse json",
			},
		)

		return
	}

	if req.CommandActionName == "" || req.DeviceName == "" {
		ctx.JSON(
			http.StatusBadRequest,
			command_dto.ResponseErr{
				Error: "could not parse json",
			},
		)

		return
	}

	creationTime := time.Now()
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			command_dto.ResponseErr{
				Error: "user not found",
			},
		)
	}

	cmd := command.Command{
		CommandActionName: req.CommandActionName,
		DeviceName:        req.DeviceName,
		Payload:           req.Payload,
		Status:            command.WAITING,
		CreatedAt:         &creationTime,
		QueuedAt:          &creationTime,
		CompletedAt:       nil,
		Priority:          req.Priority,
		RequestedBy:       currentUser.(*user.User).Username,
		ErrorDescription:  "",
	}
	id, err := handler.service.EnqueueCommand(ctx, &cmd)
	if err != nil {
		res := command_dto.ResponseErr{
			Error: "generic error",
		}
		httpStatus := http.StatusInternalServerError

		switch {
		case errors.Is(err, command.ErrDeviceNotFound):
			res.Error = err.Error()
			httpStatus = http.StatusBadRequest
		case errors.Is(err, command_action.ErrCommandActionNotFound):
			actNames, err := handler.service.GetAvailableCmdActions(ctx)
			if err != nil {
				res.Error = err.Error()
			} else {
				errStr := "command action should be one of:"
				for _, item := range actNames {
					errStr += " " + item + ","
				}
				res.Error = errStr[:len(errStr)-1]
			}

			httpStatus = http.StatusBadRequest
		}

		ctx.JSON(httpStatus, res)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		command_dto.AddNewCommadResponse{
			ID:                *id,
			CommandActionName: cmd.CommandActionName,
			Status:            cmd.Status,
		},
	)
}

// @BasePath /command

// Lists the previously enqueued commands
// @Summary Lists the previously enqueued commands applying a filter
// @Schemes
// @Description Shows a list of all enqueued commands, with their status and other information useful for debugging the command operations. The list is filtered with some parameters
// @Tags command
// @Accept json
// @Produce json
// @Param id query string false "Command ID"
// @Param deviceName query string false "Target device name"
// @Param status query string false "Status"
// @Param priority query string false "Priority"
// @Param requestedBy query string false "Requested by"
// @Success 200 {object} command_dto.ListCommandsResponse
// @Failure 400 {object} command_dto.ResponseErr
// @Failure 401 {object} command_dto.ResponseErr
// @Failure 500 {object} command_dto.ResponseErr
// @Router /command/list [get]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
func (handler *CommandHandler) ListCommands(ctx *gin.Context) {
	id := ctx.Query("id")
	deviceName := ctx.Query("deviceName")
	status := ctx.Query("status")
	commandActionName := ctx.Query("commandActioName")
	priority := ctx.Query("priority")
	requestedBy := ctx.Query("requestedBy")

	var priorityInt *uint
	if priority != "" {
		val, err := strconv.Atoi(priority)

		if err != nil || val < 0 || val >= 0xFFFF {
			ctx.JSON(http.StatusBadRequest, command_dto.ResponseErr{
				Error: "priority is not valid",
			})
			return
		}

		casted := uint(val)
		priorityInt = &casted

	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, command_dto.ResponseErr{
			Error: "id is not hex",
		})
		return
	}

	filter := command.CommandFilter{
		Id:                &objId,
		DeviceName:        deviceName,
		Status:            command.StatusFromString(status),
		CommandActionName: commandActionName,
		Priority:          priorityInt,
		RequestedBy:       requestedBy,
	}

	commands, err := handler.service.ListCommands(ctx, filter)
	if err != nil {
		resErr := command_dto.ResponseErr{
			Error: "generic error",
		}
		httpStatus := http.StatusInternalServerError

		switch {
		case errors.Is(err, command.ErrParsingResult):
			resErr.Error = err.Error()
			httpStatus = http.StatusInternalServerError
		}

		ctx.JSON(httpStatus, resErr)
		return
	}

	ctx.JSON(
		http.StatusOK,
		command_dto.ListCommandsResponse{
			Commands: commands,
		},
	)
}

func (handler *CommandHandler) DeleteCommand(ctx *gin.Context) {
	var req command_dto.DeleteCommandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			command_dto.ResponseErr{
				Error: "could not parse json",
			},
		)

		return
	}

	err := handler.service.Delete(ctx, req.ID)
	if err != nil {
		resErr := command_dto.ResponseErr{Error: "generic error"}
		httpStatus := http.StatusInternalServerError

		switch {
		case errors.Is(err, command.ErrCommandNotFound):
			resErr.Error = err.Error()
			httpStatus = http.StatusNotFound
		case errors.Is(err, command.ErrIdMalformed):
			resErr.Error = err.Error()
			httpStatus = http.StatusBadRequest
		}

		ctx.JSON(httpStatus, resErr)
		return

	}

	ctx.JSON(http.StatusOK, command_dto.DeleteCommandResponse{
		ID: req.ID,
	})
}

// @BasePath /command

// Updates a command status
// @Summary Updates an enqueued command status
// @Schemes
// @Description Lets a user update the status of an enqueue command manually (enqueued, acknowledged, completed, errored)
// @Tags command
// @Accept json
// @Produce json
// @Param command body command_dto.UpdateCommandStatusRequest true "command data"
// @Success 200 {object} command_dto.UpdateCommandStatusResponse
// @Failure 400 {object} command_dto.ResponseErr
// @Failure 401 {object} command_dto.ResponseErr
// @Failure 500 {object} command_dto.ResponseErr
// @Router /command/update/status [post]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
func (handler *CommandHandler) UpdateCommandStatus(ctx *gin.Context) {
	var req command_dto.UpdateCommandStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			command_dto.ResponseErr{
				Error: "could not parse json",
			},
		)
		return
	}

	newStatus := command.StatusFromString(req.Status)
	if newStatus == nil ||
		newStatus != &command.COMPLETED &&
			newStatus != &command.ERRORED {
		ctx.JSON(
			http.StatusBadRequest,
			command_dto.ResponseErr{
				Error: `could not parse status. 
				should be either completed or errored`,
			},
		)

		return
	}

	if newStatus == &command.ERRORED && req.ErrorDescription == "" {
		ctx.JSON(
			http.StatusBadRequest,
			command_dto.ResponseErr{
				Error: "errored status must have error description",
			},
		)

		return
	}

	errDesc := ""
	if newStatus == &command.ERRORED {
		errDesc = req.ErrorDescription
	}

	err := handler.service.UpdateStatus(ctx, req.ID, *newStatus, errDesc, "")
	if err != nil {
		httpStatus := http.StatusInternalServerError
		errRes := command_dto.ResponseErr{
			Error: "generic error",
		}

		switch {
		case errors.Is(err, command.ErrCommandNotFound):
			httpStatus = http.StatusNotFound
			errRes.Error = err.Error()
		case errors.Is(err, command.ErrIdMalformed):
			httpStatus = http.StatusBadRequest
			errRes.Error = err.Error()
		}

		ctx.JSON(httpStatus, errRes)
		return
	}

	ctx.JSON(
		http.StatusOK,
		command_dto.UpdateCommandStatusResponse{
			ID:        req.ID,
			NewStatus: req.Status,
		},
	)
}
