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

type TransactionCommitError struct {
	message string
}

func NewTransactionCommitError(msg string) TransactionCommitError {
	return TransactionCommitError{
		message: msg,
	}
}

func (err TransactionCommitError) Error() string {
	return err.message
}

type TemplateFindError struct {
	Message string
}

func (err TemplateFindError) Error() string {
	return err.Message
}

type TemplateUpdateError struct {
	Message string
}

func (err TemplateUpdateError) Error() string {
	return err.Message
}
