package common

import "errors"

var (
	// ErrNotImplemented ...
	ErrNotImplemented = errors.New("not implemented")
	// ErrInvalidParams ...
	ErrInvalidParams = errors.New("invalid parameters")
	// ErrInvalidFilePath ..
	ErrInvalidFilePath = errors.New("invalid file path")
	// ErrFileExists ...
	ErrFileExists = errors.New("file exists")
	// ErrCannotAccess ...
	ErrCannotAccess = errors.New("can't access file")
	// ErrUnexpected ...
	ErrUnexpected = errors.New("unexpected error")
)
