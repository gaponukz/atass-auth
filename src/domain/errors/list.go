package errors

import "errors"

var (
	ErrRouteNotFound               = errors.New("Route not found")
	ErrUserNotFound                = errors.New("User not found")
	ErrUserNotValid                = errors.New("User not valid")
	ErrUserAlreadyExists           = errors.New("User already exists")
	ErrPasswordResetRequestMissing = errors.New("User did not submit a password reset request")
	ErrRegisterRequestMissing      = errors.New("User did not submit a register request")
	ErrTokenEarlyToUpdate          = errors.New("Is too early to update token")
)
