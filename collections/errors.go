package collections

import (
	"errors"
	"fmt"
)

type PersistenceError struct {
	Err error
}

func (e PersistenceError) Error() string {
	return fmt.Sprintf("persistence error: %s", e.Err)
}

type NotFoundError struct {
	Message string
	Err     error
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found error: %s", e.Message)
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		Err:     errors.New(message),
		Message: message,
	}
}

func NewNotFoundErrorWithOriginalError(message string, err error) NotFoundError {
	return NotFoundError{
		Err:     err,
		Message: message,
	}
}
