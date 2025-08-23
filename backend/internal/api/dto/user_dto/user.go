package user_dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	Username string `json:"username,omitempty"`
	Error    string `json:"error,omitempty"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type UpdateUserEnableStatusRequest struct {
	Username string `json:"username" binding:"required"`
	Enabled  bool   `json:"enabled" binding:"required"`
}

type UpdateUserEnableStatusResponse struct {
	Username string `json:"username"`
	Enabled  bool   `json:"enabled"`
}

type UpdateUserEnableStatusResponseErr struct {
	Error string `json:"error"`
}
