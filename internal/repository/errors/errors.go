package errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrUniqueConstaint = errors.New("unique constraint")
)
