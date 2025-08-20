package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/dto"
	"github.com/MattBrs/OcelotMDM/internal/domain/device"
	"github.com/MattBrs/OcelotMDM/internal/domain/token"
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
		ctx.JSON(
			http.StatusBadRequest,
			dto.DeviceCreationErrResponse{
				Error: "invalid JSON",
			},
		)

		return
	}

	if req.Otp == "" || req.Type == "" || req.Architecture == "" {
		ctx.JSON(
			http.StatusBadRequest,
			dto.DeviceCreationErrResponse{
				Error: "invalid JSON",
			},
		)

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

	newCert, err := h.service.RegisterNewDevice(ctx.Request.Context(), &dev, req.Otp)
	if err != nil {
		var errRes dto.DeviceCreationErrResponse
		switch {
		case errors.Is(err, device.ErrInvalidOtp):
			errRes.Error = "the otp is no longer valid"
		case errors.Is(err, token.ErrOtpNotFound):
			errRes.Error = "the otp was not found"
		default:
			errRes.Error = "failed to register the device"
		}

		ctx.JSON(
			http.StatusInternalServerError,
			errRes,
		)

		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.DeviceCreationResponse{
			Name:     name,
			OvpnFile: string(newCert),
		},
	)
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
			dto.UpdateAddressErrResponse{Error: "Could not parse JSON"},
		)
	}

	if err := h.service.UpdateAddress(ctx, req.Name, req.IPAddress); err != nil {
		var errRes dto.UpdateAddressErrResponse
		httpStatus := http.StatusInternalServerError
		switch {
		case errors.Is(err, device.ErrDeviceNotFound):
			httpStatus = http.StatusNotFound
			errRes.Error = "device was not found"
		case errors.Is(err, device.ErrDeviceNotUpdated):
			errRes.Error = "device addr was not updated"
		default:
			errRes.Error = "generic error"
		}

		ctx.JSON(httpStatus, errRes)
		return
	}

	ctx.JSON(
		http.StatusOK,
		dto.UpdateAddressResponse{
			DeviceName: req.Name,
			IpAddress:  req.IPAddress,
		},
	)
}
