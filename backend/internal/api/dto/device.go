package dto

type DeviceCreationRequest struct {
	Otp          string `json:"otp"`
	Type         string `json:"type"`
	Architecture string `json:"architecture"`
}

type DeviceCreationResponse struct {
	// TODO: add needed client certificates and configurations in the response
	Name  string `json:"name"`
	Error string `json:"error,omitempty"`
}

type UpdateAddressRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type UpdateAddressResponse struct {
	Error string `json:"error,omitempty"`
}
