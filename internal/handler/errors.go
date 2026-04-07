package handler

import "errors"

var (
	// request error
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrInvalidRequestBody = errors.New("invalid request body")

	// server side error
	ErrInternalServerErr = errors.New("internal server error")

	// auth error
	ErrNotAuthorized = errors.New("not authorized")
)
