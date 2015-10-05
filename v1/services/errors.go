package services

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/cf"
)

type MissingKindOrClientError struct {
	Err error
}

func (e MissingKindOrClientError) Error() string {
	return e.Err.Error()
}

type CriticalKindError struct {
	Err error
}

func (e CriticalKindError) Error() string {
	return e.Err.Error()
}

type ClientMissingError struct {
	Err error
}

func (e ClientMissingError) Error() string {
	return e.Err.Error()
}

type KindMissingError struct {
	Err error
}

func (e KindMissingError) Error() string {
	return e.Err.Error()
}

func CCErrorFor(err error) error {
	if failure, ok := err.(cf.Failure); ok {
		if failure.Code == http.StatusNotFound {
			return CCNotFoundError{err}
		}
		return CCDownError{err}
	}
	return err
}

type CCNotFoundError struct {
	Err error
}

func (e CCNotFoundError) Error() string {
	return e.Err.Error()
}

type CCDownError struct {
	Err error
}

func (e CCDownError) Error() string {
	return e.Err.Error()
}

type DefaultScopeError struct{}

func (d DefaultScopeError) Error() string {
	return "You cannot send a notification to a default scope"
}
