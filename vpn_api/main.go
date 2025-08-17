package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type NewCertRequest struct {
	DeviceName string `json:"device_name"`
}

func generateClientCert(deviceName string) ([]byte, error) {
	pwd := os.Getenv("CA_PWD")
	if pwd == "" {
		return nil, errors.New("password was not provided")
	}

	bashScript := fmt.Sprintf(`
cd /etc/openvpn
export EASYRSA_BATCH=1
export EASYRSA_REQ_CN="%s"
export EASYRSA_CERT_EXPIRE=3650
export EASYRSA_NO_PASS=1
printf "%s\nyes\n" | easyrsa --batch build-client-full %s nopass
    `, deviceName, pwd, deviceName)

	certCmd := exec.Command("bash", "-c", bashScript)

	_, err := certCmd.CombinedOutput()
	if err != nil {
		return nil, errors.New("generate error")
	}

	ovpnCmd := exec.Command("ovpn_getclient", deviceName)
	ovpnCmd.Dir = "/etc/openvpn"

	output, err := ovpnCmd.Output()
	if err != nil {
		return nil, errors.New("retrieval error")
	}

	return output, nil
}

func CreateNewCert(ctx *gin.Context) {
	var req NewCertRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "500"})
		return
	}

	if req.DeviceName == "" {
		fmt.Println("deviceName is not valid")
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "device name is not valid"},
		)
		return
	}

	cert, err := generateClientCert(req.DeviceName)
	if err != nil {
		// printed error is either invalid pwd, gen error or retrieval err
		errStr := fmt.Sprintf(
			"internal server error on certificate generation: %s",
			err.Error(),
		)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": errStr},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"ovpn_file": cert},
	)
}

func main() {
	router := gin.Default()

	router.POST("/vpn/client/new", CreateNewCert)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
