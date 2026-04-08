package handler

import "errors"

var (
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInternalServerErr  = errors.New("internal server error")
	ErrUnauthorized       = errors.New("not authorized")
)
