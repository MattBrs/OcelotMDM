package token

import "errors"

var (
	ErrOtpNotFound = errors.New("OTP was not found")
	ErrOtpExpired  = errors.New("OTP is expired")
)
