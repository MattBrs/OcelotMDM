package command

import "errors"

var (
	ErrPayloadRequired = errors.New("payload is required for this command action")
	ErrCommandNotFound = errors.New("command was not found")
	ErrIdMalformed     = errors.New("inserted id is not in the correct format")
	ErrDeviceNotFound  = errors.New("device was not found")
	ErrParsingResult   = errors.New("internal server error when parsng db query")
	ErrUpdateCommand   = errors.New("an error occurred while updating the command")
)
