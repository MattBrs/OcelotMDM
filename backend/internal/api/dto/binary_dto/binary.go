package binary_dto

type AddBinaryRequest struct {
	Name         string `json:"name" binding:"required"`
	Architecture string `json:"architecture" binding:"required"`
	Version      string `json:"version" binding:"required"`
	Data         string `json:"data" binding:"required"`
}

type AddBinaryResponse struct {
	Name string `json:"name" binding:"required"`
}

type GetBinaryRequest struct {
}

type GetBinaryResponse struct {
	BinaryName string `json:"binary_name"`
	Data       []byte `json:"binary_data"`
	Version    string `json:"binary_version"`
}

type ResponseErr struct {
	Error string `json:"error"`
}
