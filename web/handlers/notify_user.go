package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
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

    results := []map[string]string{
        map[string]string{
            "status": "delivered",
        },
    }
    response, err := json.Marshal(results)
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}

func (handler NotifyUser) sendEmailTo(recipient string) {
    source, err := template.New("emailTemplate").Parse(`From: {{.From}}\nTo: {{.To}}\n`)
    if err != nil {
        panic(err)
    }
    sender := "no-reply@notifications.example.com"

    context := map[string]string{
        "From": sender,
        "To":   recipient,
    }

    buffer := bytes.NewBuffer([]byte{})
    source.Execute(buffer, context)
    message := buffer.Bytes()

    handler.logger.Print(string(message))

    env := config.NewEnvironment()
    auth := smtp.PlainAuth("", env.SMTPUser, env.SMTPPass, env.SMTPHost)
    err = smtp.SendMail(fmt.Sprintf("%s:%s", env.SMTPHost, env.SMTPPort), auth, sender, []string{recipient}, message)
    if err != nil {
        panic(err)
    }
}
