package dto

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Username string `json:"username,omitempty"`
	Error    string `json:"error,omitempty"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type UpdateUserEnableStatusRequest struct {
	Username string `json:"username"`
	Enabled  bool   `json:"enabled"`
}

type UpdateUserEnableStatusResponse struct {
	Username string `json:"username"`
	Enabled  bool   `json:"enabled"`
}

type UpdateUserEnableStatusResponseErr struct {
	Error string `json:"error"`
}
