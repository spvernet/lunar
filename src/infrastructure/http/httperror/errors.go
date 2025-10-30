package httperror

import "errors"

var (
	ErrChannelNotFound = errors.New("channel not found")
	ErrInvalidSort     = errors.New("ort by should be channel or speed or updated_at")
	ErrInvalidOrder    = errors.New("order should be asc or desc")
)
