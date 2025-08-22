package command_action

import "errors"

var (
	ErrCommandActionNotFound = errors.New("command type was not found")
	ErrNameEmpty             = errors.New("command name is empty")
	ErrDescriptionEmpty      = errors.New("command description is empty")
	ErrReqOnlineEmpty        = errors.New("command online required is not specified")
	ErrParsingCmd            = errors.New("error while parsing command")
)
