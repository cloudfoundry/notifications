package postal

import (
    "net/http"
    "net/url"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UAA string

type UAAScopesError string

func (err UAAScopesError) Error() string {
    return string(err)
}

type CCDownError string

func (err CCDownError) Error() string {
    return string(err)
}

type UAADownError string

func (err UAADownError) Error() string {
    return string(err)
}

type UAAUserNotFoundError string

func (err UAAUserNotFoundError) Error() string {
    return string(err)
}

type UAAGenericError string

func (err UAAGenericError) Error() string {
    return string(err)
}

type TemplateLoadError string

func (err TemplateLoadError) Error() string {
    return string(err)
}

type CCNotFoundError string

func (err CCNotFoundError) Error() string {
    return "CloudController Error: " + string(err)
}

func UAAErrorFor(err error) error {
    switch err.(type) {
    case *url.Error:
        return UAADownError("UAA is unavailable")
    case uaa.Failure:
        failure := err.(uaa.Failure)

        if failure.Code() == http.StatusNotFound {
            if strings.Contains(failure.Message(), "Requested route") {
                return UAADownError("UAA is unavailable")
            } else {
                return UAAGenericError("UAA Unknown 404 error message: " + failure.Message())
            }
        }

        return UAADownError(failure.Message())
    default:
        return UAAGenericError("UAA Unknown Error: " + err.Error())
    }
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

type CriticalNotificationError struct {
    kindID string
}

func NewCriticalNotificationError(kindID string) CriticalNotificationError {
    return CriticalNotificationError{kindID: kindID}
}

func (err CriticalNotificationError) Error() string {
    return "Insufficient privileges to send notification " + err.kindID
}
