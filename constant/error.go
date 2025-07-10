package constant

import "errors"

var (
	ErrInternalServerError = "Internal Server Error"
	ErrBadRequest          = "Bad Request"
	ErrNotFound            = "Not Found"
	ErrUnauthorized        = "Unauthorized"
)

var (
	ErrJwtSigned = errors.New("JWT signed error")
)
