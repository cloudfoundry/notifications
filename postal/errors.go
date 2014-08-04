package postal

import (
    "net/http"
    "net/url"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

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

        return UAADownError("UAA is unavailable")
    default:
        return UAAGenericError("UAA Unknown Error: " + err.Error())
    }
}
