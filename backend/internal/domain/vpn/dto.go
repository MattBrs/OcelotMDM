package vpn

type CreateClientRequest struct {
	DeviceName string `json:"device_name"`
}

type CreateClientResponse struct {
	OvpnFile []byte `json:"ovpn_file"`
}
