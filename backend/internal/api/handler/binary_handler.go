package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MattBrs/OcelotMDM/internal/api/dto/binary_dto"
	"github.com/MattBrs/OcelotMDM/internal/domain/binary"
	"github.com/MattBrs/OcelotMDM/internal/domain/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BinaryHandler struct {
	service *binary.Service
}

func NewBinaryHandler(service *binary.Service) *BinaryHandler {
	return &BinaryHandler{
		service,
	}
}

func (h *BinaryHandler) AddNewBinary(ctx *gin.Context) {
	var req binary_dto.AddBinaryRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, binary_dto.ResponseErr{
			Error: "could not parse JSON",
		})
		return
	}

	bin := binary.Binary{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Architecture: req.Architecture,
		Version:      req.Version,
	}

	err = h.service.AddBinary(ctx, bin, []byte(req.Data))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, binary_dto.ResponseErr{
			Error: "could not add binary to file_repo",
		})
		return
	}

	ctx.JSON(http.StatusCreated, binary_dto.AddBinaryResponse{
		Name: req.Name,
	})
}

func (h *BinaryHandler) GetBinary(ctx *gin.Context) {
	otp := ctx.Query("token")
	binaryName := ctx.Query("name")

	if otp == "" {
		ctx.JSON(http.StatusBadRequest, binary_dto.ResponseErr{
			Error: "'token' query param is missing or invalid",
		})
		return
	}

	if binaryName == "" {
		ctx.JSON(http.StatusBadRequest, binary_dto.ResponseErr{
			Error: "'name' query param is missing or invalid",
		})
		return
	}

	data, version, err := h.service.GetBinary(ctx, binaryName, otp)
	if err != nil {
		var response binary_dto.ResponseErr
		var httpStatus int

		switch {
		case errors.Is(err, token.ErrOtpNotFound):
			response.Error = "otp token was not found"
			httpStatus = http.StatusUnauthorized
		case errors.Is(err, token.ErrOtpExpired):
			response.Error = "otp token is expired"
			httpStatus = http.StatusUnauthorized
		default:
			fmt.Println("error on getBinary: ", err.Error())
			response.Error = "generic error"
			httpStatus = http.StatusInternalServerError

		}
		ctx.JSON(httpStatus, response)
		return
	}

	response := binary_dto.GetBinaryResponse{
		BinaryName: binaryName,
		Data:       data,
		Version:    *version,
	}

	ctx.JSON(http.StatusOK, response)
}
