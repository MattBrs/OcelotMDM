package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user was not found")
	ErrCouldNotHashPwd   = errors.New("error hashing pwd")
	ErrFailedToConvertID = errors.New("failed to convert objectID to primitive")
	ErrUserNotUpdated    = errors.New("error while updating user")
	ErrUsernameNotValid  = errors.New("selected username is not valid")
	ErrPasswordNotValid  = errors.New("inserted password is not valid")
	ErrUsernameTaken     = errors.New("inserted username is already taken")
	ErrCreatingUser      = errors.New("an error occurred while creating the user")
	ErrTokenGeneration   = errors.New("error while generating token")
)
