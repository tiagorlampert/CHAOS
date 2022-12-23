package internal

import "errors"

var (
	ErrNoDatabaseProvided       = errors.New("no database provided")
	ErrClientConnectionNotFound = errors.New("no active client connection found")
	ErrInvalidServerAddress     = errors.New("the server address provided is invalid")
	ErrInvalidServerPort        = errors.New("the server port provided is invalid")
)
