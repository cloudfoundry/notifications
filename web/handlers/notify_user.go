package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/dgrijalva/jwt-go"
    uuid "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const emailTemplate = `The following "{{.KindDescription}}" notification was sent to you directly by the "{{.SourceDescription}}" component of Cloud Foundry:

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
    params := NewNotifyUserParams(req.Body)

    if valid := params.Validate(); !valid {
        handler.Error(w, 422, params.Errors)
        return
    }

    token, err := handler.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    handler.uaaClient.SetToken(token.Access)

    user, ok := handler.loadUser(w, req)
    if !ok {
        return
    }

    status := handler.sendMailToUser(user, params, req)
    response := handler.generateResponse(status)

    w.WriteHeader(http.StatusOK)
    w.Write(response)
}

func (handler NotifyUser) generateResponse(status string) []byte {
    results := []map[string]string{
        map[string]string{
            "status": status,
        },
    }
    response, err := json.Marshal(results)
    if err != nil {
        panic(err)
    }
    return response
}

func (handler NotifyUser) loadUser(w http.ResponseWriter, req *http.Request) (uaa.User, bool) {
    userID := strings.TrimPrefix(req.URL.Path, "/users/")
    user, err := handler.uaaClient.UserByID(userID)
    if err != nil {
        switch err.(type) {
        case *url.Error:
            w.WriteHeader(http.StatusBadGateway)
        case uaa.Failure:
            w.WriteHeader(http.StatusGone)
        default:
            w.WriteHeader(http.StatusInternalServerError)
        }
        return uaa.User{}, false
    }
    return user, true
}

func (handler NotifyUser) sendMailToUser(user uaa.User, params NotifyUserParams, req *http.Request) string {
    env := config.NewEnvironment()
    var status string
    var message mail.Message
    if len(user.Emails) > 0 {
        guid, err := handler.guidGenerator()
        if err != nil {
            panic(err)
        }

        rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
        token, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
            return []byte(config.UAAPublicKey), nil
        })

        context := MessageContext{
            From:              env.Sender,
            To:                user.Emails[0],
            Subject:           params.Subject,
            Text:              params.Text,
            Template:          emailTemplate,
            KindDescription:   params.KindDescription,
            SourceDescription: params.SourceDescription,
            ClientID:          token.Claims["client_id"].(string),
            MessageID:         guid.String(),
        }
        handler.logger.Printf("Sending email to %s", context.To)
        status, message, err = SendMail(handler.mailClient, context)
        if err != nil {
            panic(err)
        }

        handler.logger.Print(message.Data())
    }
    return status
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
