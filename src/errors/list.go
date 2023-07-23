package errors

import "errors"

var (
	ErrUserNotFound                = errors.New("User not found")
	ErrUserAlreadyExists           = errors.New("User already exists")
	ErrPasswordResetRequestMissing = errors.New("User did not submit a password reset request")
	ErrRegisterRequestMissing      = errors.New("User did not submit a register request")
	ErrRouteNotFound               = errors.New("Can not find route")
)
