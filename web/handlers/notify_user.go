package handlers

import (
    "encoding/json"
    "log"
    "net/http"
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
        Error(w, 422, params.Errors)
        return
    }

    token, err := handler.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    handler.uaaClient.SetToken(token.Access)

    userGUID := strings.TrimPrefix(req.URL.Path, "/users/")
    user, ok := loadUser(w, userGUID, handler.uaaClient)
    if !ok {
        return
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    var status string
    if len(user.Emails) > 0 {
        env := config.NewEnvironment()
        context := handler.buildContext(user, params, env, clientToken.Claims["client_id"].(string))
        status = sendMailToUser(context, handler.logger, handler.mailClient)
    } else {
        status = "User had no email address"
    }

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

func (handler NotifyUser) buildContext(user uaa.User, params NotifyUserParams, env config.Environment, clientID string) MessageContext {
    guid, err := handler.guidGenerator()
    if err != nil {
        panic(err)
    }

    return MessageContext{
        From:              env.Sender,
        To:                user.Emails[0],
        Subject:           params.Subject,
        Text:              params.Text,
        Template:          emailTemplate,
        KindDescription:   params.KindDescription,
        SourceDescription: params.SourceDescription,
        ClientID:          clientID,
        MessageID:         guid.String(),
    }
}
