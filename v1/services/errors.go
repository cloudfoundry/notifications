package services

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/cf"
)

type MissingKindOrClientError string

func (err MissingKindOrClientError) Error() string {
	return string(err)
}

type CriticalKindError string

func (err CriticalKindError) Error() string {
	return string(err)
}

type ClientMissingError string

func (err ClientMissingError) Error() string {
	return string(err)
}

type KindMissingError string

func (err KindMissingError) Error() string {
	return string(err)
}

type TemplateAssignmentError string

func (err TemplateAssignmentError) Error() string {
	return string(err)
}

func CCErrorFor(err error) error {
	if failure, ok := err.(cf.Failure); ok {
		if failure.Code == http.StatusNotFound {
			return CCNotFoundError(err.Error())
		}
		return CCDownError(err.Error())
	}
	return err
}

type CCNotFoundError string

func (err CCNotFoundError) Error() string {
	return "CloudController Error: " + string(err)
}

type CCDownError string

func (err CCDownError) Error() string {
	return string(err)
}
