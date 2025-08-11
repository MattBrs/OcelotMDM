package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/dto"
	"github.com/MattBrs/OcelotMDM/internal/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			dto.CreateUserResponse{Error: "Invalid JSON"},
		)
		return
	}

	newUser := user.User{
		Username:  req.Username,
		Password:  req.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UpdatedBy: primitive.NewObjectID(),
		Enabled:   false,
	}

	res := dto.CreateUserResponse{}
	err := h.service.CreateNewUser(ctx.Request.Context(), &newUser)
	if err != nil {
		fmt.Println(err.Error())
		switch {
		case errors.Is(err, user.ErrUsernameTaken):
			res.Error = "username already taken"
		case errors.Is(err, user.ErrPasswordNotValid):
			res.Error = "passwords does not follow security guidelines"
		case errors.Is(err, user.ErrUsernameNotValid):
			res.Error = "username is not valid"
		case errors.Is(err, user.ErrFailedToConvertID):
			res.Error = "failed to convert id"
		default:
			res.Error = "generic error"
		}

		ctx.JSON(
			http.StatusInternalServerError,
			res,
		)
		return
	}

	res.Username = newUser.Username
	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req dto.LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			dto.LoginUserResponse{
				Error: "could not parse JSON",
			},
		)
		return
	}

	var res dto.LoginUserResponse
	token, err := h.service.LoginUser(ctx, req.Username, req.Password)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			status = http.StatusUnauthorized
			res.Error = "username or password are not correct"
		case errors.Is(err, user.ErrPasswordNotValid):
			status = http.StatusUnauthorized
			res.Error = "username or password are not correct"
		case errors.Is(err, user.ErrTokenGeneration):
			res.Error = "there was an error while generating the token"
		case errors.Is(err, user.ErrUserNotAuthorized):
			res.Error = "user is not enabled. contact an administrator"
		default:
			res.Error = "generic error"
		}

		ctx.JSON(status, res)
		return
	}
	res.Token = *token
	ctx.JSON(http.StatusOK, res)
}
