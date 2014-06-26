package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type NotifyUser struct {
    logger *log.Logger
}

func NewNotifyUser(logger *log.Logger) NotifyUser {
    return NotifyUser{
        logger: logger,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    env := config.NewEnvironment()
    uaaConfig := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, req.Header.Get("Authorization"))
    user, err := uaa.UserByID(uaaConfig, strings.Trim(req.URL.Path, "/users/"))
    if err != nil {
        panic(err)
    }

    for _, email := range user.Emails {
        handler.logger.Printf("Sending email to %s", email)
    }
}
