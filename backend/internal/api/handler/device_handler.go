package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/dto/device_dto"
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

// @BasePath /devices

// Enroll a new device
// @Summary enroll a new device to the network
// @Schemes
// @Description enroll device
// @Tags devices
// @Accept json
// @Produce json
// @Param device body device_dto.DeviceCreationRequest true "Device Data"
// @Success 200 {object} device_dto.DeviceCreationResponse
// @Router /devices [post]
func (h *DeviceHandler) AddNewDevice(ctx *gin.Context) {
	var req device_dto.DeviceCreationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			device_dto.DeviceCreationErrResponse{
				Error: "invalid JSON",
			},
		)

		return
	}

	if req.Otp == "" || req.Type == "" || req.Architecture == "" {
		ctx.JSON(
			http.StatusBadRequest,
			device_dto.DeviceCreationErrResponse{
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
		var errRes device_dto.DeviceCreationErrResponse
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
		device_dto.DeviceCreationResponse{
			Name:     name,
			OvpnFile: string(newCert),
		},
	)
}

// @BasePath /devices

// List devices
// @Summary show a filtered list of the devices
// @Schemes
// @Description show a filtered list of the devices previously enrolled
// @Tags devices
// @Accept json
// @Produce json
// @Param id query string false "Device ID"
// @Param status query string false "Device Status"
// @Param name query string false "Device Name"
// @Param architecture query string false "Device Architecture"
// @Success 200 {object} []device.Device
// @Router /devices [get]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
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

// @BasePath /devices

// Update ip address
// @Summary update device address
// @Schemes
// @Description update known device ip address
// @Tags devices
// @Accept json
// @Produce json
// @Param device body device_dto.UpdateAddressRequest true "Device New address"
// @Success 200 {object} device_dto.DeviceCreationResponse
// @Router /devices/updateAddress [post]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
func (h *DeviceHandler) UpdateDeviceAddress(ctx *gin.Context) {
	var req device_dto.UpdateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			device_dto.UpdateAddressErrResponse{Error: "Could not parse JSON"},
		)
	}

	if err := h.service.UpdateAddress(ctx, req.Name, req.IPAddress); err != nil {
		var errRes device_dto.UpdateAddressErrResponse
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
		device_dto.UpdateAddressResponse{
			DeviceName: req.Name,
			IpAddress:  req.IPAddress,
		},
	)
}
