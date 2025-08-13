package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/dto"
	"github.com/MattBrs/OcelotMDM/internal/device"
	"github.com/MattBrs/OcelotMDM/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/goombaio/namegenerator"
)

type DeviceHandler struct {
	service   *device.Service
	generator namegenerator.Generator
}

func NewDeviceHandler(service *device.Service) *DeviceHandler {
	seed := time.Now().UTC().UnixNano()
	return &DeviceHandler{service, namegenerator.NewNameGenerator(seed)}
}

func (h *DeviceHandler) AddNewDevice(ctx *gin.Context) {
	var req dto.DeviceCreationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	name := h.generator.Generate()
	dev := device.Device{
		Name:         name,
		Type:         req.Type,
		IPAddress:    "Unk",
		Status:       "Unk",
		LastSeen:     time.Now().Unix(),
		Tags:         []string{},
		Architecture: req.Architecture,
	}

	err := h.service.RegisterNewDevice(ctx.Request.Context(), &dev, req.Otp)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrInvalidOtp):
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "the otp is no longer valid"},
			)
		case errors.Is(err, token.ErrOtpNotFound):
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "the otp was not found"},
			)
		default:
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "failed to register the device"},
			)
		}

		return
	}

	res := dto.DeviceCreationResponse{
		Name: name,
	}
	ctx.JSON(http.StatusCreated, res)
}

func (h *DeviceHandler) ListDevices(ctx *gin.Context) {
	id := ctx.Query("id")
	status := ctx.Query("status")
	name := ctx.Query("name")
	architecture := ctx.Query("architecture")

	deviceFilter := device.DeviceFilter{
		Id:           id,
		Status:       status,
		Name:         name,
		Architecture: architecture,
	}

	devices, err := h.service.ListDevices(ctx.Request.Context(), deviceFilter)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Could not fetch devices"},
		)
		return
	}

	ctx.JSON(http.StatusFound, devices)
}

func (h *DeviceHandler) UpdateDeviceAddress(ctx *gin.Context) {
	var req dto.UpdateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			dto.UpdateAddressResponse{Error: "Could not parse JSON"},
		)
	}

	if err := h.service.UpdateAddress(ctx, req.Name, req.IPAddress); err != nil {
		switch {
		case errors.Is(err, device.ErrDeviceNotFound):
			ctx.JSON(
				http.StatusInternalServerError,
				dto.DeviceCreationResponse{Error: "device was not found"},
			)
		case errors.Is(err, device.ErrDeviceNotUpdated):
			ctx.JSON(
				http.StatusInternalServerError,
				dto.DeviceCreationResponse{Error: "device addr was not updated"},
			)
		default:
			ctx.JSON(
				http.StatusInternalServerError,
				dto.DeviceCreationResponse{Error: "generic error"},
			)
		}

		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.UpdateAddressResponse{},
	)
}
