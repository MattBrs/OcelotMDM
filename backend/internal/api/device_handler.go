package api

import (
	"net/http"

	"github.com/MattBrs/OcelotMDM/internal/device"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	service *device.Service
}

func newDeviceHandler(service *device.Service) *DeviceHandler {
	return &DeviceHandler{service}
}

func (h *DeviceHandler) addNewDevice(ctx *gin.Context) {
	var req device.Device
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := h.service.RegisterNewDevice(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add the device"})
		return
	}

	ctx.JSON(http.StatusCreated, req)
}
