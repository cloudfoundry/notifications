package webutil

import "fmt"

type ParseError struct{}

func (err ParseError) Error() string {
	return "Request body could not be parsed"
}

type SchemaError struct {
	Err error
}

func (e SchemaError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type MissingUserTokenError struct {
	Err error
}

func (e MissingUserTokenError) Error() string {
	return e.Err.Error()
}

type TemplateCreateError struct{}

func (err TemplateCreateError) Error() string {
	return "Failed to create Template in the database"
}

type UAAScopesError struct {
	Err error
}

func (e UAAScopesError) Error() string {
	return e.Err.Error()
}

type CriticalNotificationError struct {
	Err error
}

func NewCriticalNotificationError(kindID string) CriticalNotificationError {
	return CriticalNotificationError{fmt.Errorf("Insufficient privileges to send notification %s", kindID)}
}

func (e CriticalNotificationError) Error() string {
	return e.Err.Error()
}
