package device

import "errors"

var (
	ErrInvalidOtp       = errors.New("invalid OTP")
	ErrEmptyName        = errors.New("name is empty")
	ErrEmptyType        = errors.New("type is empty")
	ErrDeviceNotFound   = errors.New("device not found with filter")
	ErrDeviceNotUpdated = errors.New("device was not updated")
)
