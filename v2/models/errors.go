package models

import "fmt"

type RecordNotFoundError struct {
	Err error
}

func NewRecordNotFoundError(format string, args ...interface{}) RecordNotFoundError {
	return RecordNotFoundError{fmt.Errorf(format, args...)}
}

func (err RecordNotFoundError) Error() string {
	return err.Err.Error()
}

type DuplicateRecordError struct {
	Err error
}

func (err DuplicateRecordError) Error() string {
	return err.Err.Error()
}
