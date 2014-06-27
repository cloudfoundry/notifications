package handlers

import (
    "bytes"
    "log"
    "net/http"
    "strings"
    "text/template"

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
    uaaConfig := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "))
    user, err := uaa.UserByID(uaaConfig, strings.TrimPrefix(req.URL.Path, "/users/"))
    if err != nil {
        panic(err)
    }

    if len(user.Emails) > 0 {
        address := user.Emails[0]
        handler.logger.Printf("Sending email to %s", address)

        handler.sendEmailTo(address)
    }
}

func (handler NotifyUser) sendEmailTo(recipient string) {
    source, err := template.New("emailTemplate").Parse(`From: {{.From}}\nTo: {{.To}}\n`)
    if err != nil {
        panic(err)
    }

    context := map[string]string{
        "From": "notifications@cf-app.com",
        "To":   recipient,
    }

    buffer := bytes.NewBuffer([]byte{})
    source.Execute(buffer, context)

    handler.logger.Print(buffer.String())
}
