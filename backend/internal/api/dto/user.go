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
