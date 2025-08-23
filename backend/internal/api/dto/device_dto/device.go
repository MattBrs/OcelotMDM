package device_dto

type DeviceCreationRequest struct {
	Otp          string `json:"otp" binding:"required"`
	Type         string `json:"type" binding:"required"`
	Architecture string `json:"architecture" binding:"required"`
}

type DeviceCreationResponse struct {
	// TODO: add needed client certificates and configurations in the response
	Name     string `json:"name" binding:"required"`
	OvpnFile string `json:"ovpn_file" binding:"required"`
}

type DeviceCreationErrResponse struct {
	Error string `json:"error,omitempty"`
}

type UpdateAddressRequest struct {
	Name      string `json:"name" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
}

type UpdateAddressErrResponse struct {
	Error string `json:"error"`
}

type UpdateAddressResponse struct {
	DeviceName string `json:"device_name"`
	IpAddress  string `json:"ip_address"`
}
