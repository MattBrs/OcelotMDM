package dto

type DeviceCreationRequest struct {
	Otp  string `json:"otp"`
	Type string `json:"type"`
}

type DeviceCreationResponse struct {
	// TODO: add needed client certificates and configurations in the response
	Name string ``
}
