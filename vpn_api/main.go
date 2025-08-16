package vpnapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewCertRequest struct {
	DeviceName string `bson:"device_name"`
}

func CreateNewCert(ctx *gin.Context) {
	var req NewCertRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": "200"})
}

func Main() {
	router := gin.Default()

	router.POST("/new", CreateNewCert)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
