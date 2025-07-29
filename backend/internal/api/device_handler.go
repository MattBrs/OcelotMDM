package api

import (
	"net/http"
	"strconv"

	"github.com/MattBrs/OcelotMDM/internal/device"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	service *device.Service
}

func NewDeviceHandler(service *device.Service) *DeviceHandler {
	return &DeviceHandler{service}
}

func (h *DeviceHandler) AddNewDevice(ctx *gin.Context) {
	var req device.Device
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := h.service.RegisterNewDevice(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add the device"})
		return
	}

	ctx.JSON(http.StatusCreated, req)
}

func (h *DeviceHandler) ListDevices(ctx *gin.Context) {
	id := ctx.Query("id")
	status := ctx.Query("status")
	name := ctx.Query("name")
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 0
	}

	deviceFilter := device.DeviceFilter{
		Id:     id,
		Status: status,
		Name:   name,
		Limit:  limit,
	}

	devices, err := h.service.ListDevices(ctx.Request.Context(), deviceFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add the device"})
	}

	ctx.JSON(http.StatusFound, gin.H{"error": "TODO"})
}
