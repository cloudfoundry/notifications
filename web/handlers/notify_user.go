package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
    "text/template"

    "github.com/cloudfoundry-incubator/notifications/mail"
    uuid "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const emailBody = `The following "{{.KindDescription}}" notification was sent to you directly by the "{{.SourceDescription}}" component of Cloud Foundry:

{{.Text}}`

type GUIDGenerationFunc func() (*uuid.UUID, error)

type NotifyUser struct {
    logger        *log.Logger
    mailClient    mail.ClientInterface
    uaaClient     uaa.UAAInterface
    guidGenerator GUIDGenerationFunc
}

func NewNotifyUser(logger *log.Logger, mailClient mail.ClientInterface, uaaClient uaa.UAAInterface, guidGenerator GUIDGenerationFunc) NotifyUser {
    return NotifyUser{
        logger:        logger,
        mailClient:    mailClient,
        uaaClient:     uaaClient,
        guidGenerator: guidGenerator,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params := NewNotifyUserParams(req)
    params.ParseRequestPath()
    params.ParseEnvironmentVariables()

    params.ParseRequestBody()
    if valid := params.ValidateRequestBody(); !valid {
        handler.Error(w, 422, params.Errors)
        return
    }

    params.ParseAuthorizationToken()
    if valid := params.ValidateAuthorizationToken(); !valid {
        handler.Error(w, http.StatusForbidden, params.Errors)
        return
    }

    token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    handler.uaaClient.SetToken(token)
    user, err := handler.uaaClient.UserByID(params.UserID)
    if err != nil {
        switch err.(type) {
        case *url.Error:
            w.WriteHeader(http.StatusBadGateway)
        case uaa.Failure:
            w.WriteHeader(http.StatusGone)
        default:
            w.WriteHeader(http.StatusInternalServerError)
        }
        return
    }

    var status string
    if len(user.Emails) > 0 {
        params.To = user.Emails[0]
        handler.logger.Printf("Sending email to %s", params.To)
        status = handler.SendEmail(params)
    }

    results := []map[string]string{
        map[string]string{
            "status": status,
        },
    }
    response, err := json.Marshal(results)
    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(response)
}

func (handler NotifyUser) SendEmail(context NotifyUserParams) string {
    source, err := template.New("emailBody").Parse(emailBody)
    buffer := bytes.NewBuffer([]byte{})
    source.Execute(buffer, context)

    guid, err := handler.guidGenerator()
    if err != nil {
        panic(err)
    }

    msg := mail.Message{
        From:    context.From,
        To:      context.To,
        Subject: fmt.Sprintf("CF Notification: %s", context.Subject),
        Body:    buffer.String(),
        Headers: []string{
            fmt.Sprintf("X-CF-Client-ID: %s", context.ClientID),
            fmt.Sprintf("X-CF-Notification-ID: %s", guid.String()),
        },
    }

    handler.logger.Print(msg.Data())

    err = handler.mailClient.Connect()
    if err != nil {
        return "unavailable"
    }

    err = handler.mailClient.Send(msg)
    if err != nil {
        return "failed"
    }

    return "delivered"
}

func (handler NotifyUser) Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}
