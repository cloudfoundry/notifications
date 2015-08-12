package models

import "fmt"

func NewRecordNotFoundError(format string, arguments ...interface{}) RecordNotFoundError {
	return RecordNotFoundError(fmt.Sprintf(format, arguments...))
}

type RecordNotFoundError string

func (err RecordNotFoundError) Error() string {
	message := "Record Not Found"
	if err != "" {
		message = fmt.Sprintf("%s: %s", message, string(err))
	}

	return message
}

type DuplicateRecordError struct{}

func (err DuplicateRecordError) Error() string {
	return "Duplicate Record"
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
