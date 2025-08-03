package device

import "errors"

var (
	ErrInvalidOtp = errors.New("invalid OTP")
	ErrEmptyName  = errors.New("name is empty")
	ErrEmptyType  = errors.New("type is empty")
)
