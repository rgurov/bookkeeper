package errors

import "errors"

var (
	ErrNotFound = errors.New("err not found")
	ErrDataBusy = errors.New("data is busy")
)
