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

func generateClientCert(deviceName string) (string, error) {
	pwd := os.Getenv("CA_PWD")
	if pwd == "" {
		return "", errors.New("password was not provided")
	}

	fmt.Println("pwd:", pwd)
	fmt.Println("deviceName:", deviceName)

	bashScript := fmt.Sprintf(`
cd /etc/openvpn
export EASYRSA_BATCH=1
export EASYRSA_REQ_CN="%s"
export EASYRSA_CERT_EXPIRE=3650
export EASYRSA_NO_PASS=1
printf "%s\nyes\n" | easyrsa --batch build-client-full %s nopass
    `, deviceName, pwd, deviceName)

	certCmd := exec.Command("bash", "-c", bashScript)

	output, err := certCmd.CombinedOutput()
	if err != nil {
		fmt.Println("gen conf error: ", string(output))
		return "", err
	}

	ovpnCmd := exec.Command("ovpn_getclient", deviceName)
	ovpnCmd.Dir = "/etc/openvpn"

	output, err = ovpnCmd.Output()
	if err != nil {
		fmt.Println("conf pull error: ", string(output))
		return "", err
	}

	return string(output), nil
}

func CreateNewCert(ctx *gin.Context) {
	var req NewCertRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Println("err: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500"})
		return
	}

	if req.DeviceName == "" {
		fmt.Println("deviceName is void")
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500"})
		return
	}

	cert, err := generateClientCert(req.DeviceName)
	if err != nil {
		fmt.Println("err generate cert: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500"})
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"code": "200", "ovpn_file": cert},
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
