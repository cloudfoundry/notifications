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
    helper        NotifyHelper
}

func NewNotifyUser(logger *log.Logger, mailClient mail.ClientInterface, uaaClient uaa.UAAInterface, guidGenerator GUIDGenerationFunc) NotifyUser {
    return NotifyUser{
        logger:        logger,
        mailClient:    mailClient,
        uaaClient:     uaaClient,
        guidGenerator: guidGenerator,
        helper:        NotifyHelper{},
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params := NewNotifyParams(req.Body)

    if valid := params.Validate(); !valid {
        handler.helper.Error(w, 422, params.Errors)
        return
    }

    token, err := handler.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    handler.uaaClient.SetToken(token.Access)

    userGUID := strings.TrimPrefix(req.URL.Path, "/users/")
    user, ok := handler.helper.LoadUser(w, userGUID, handler.uaaClient)
    if !ok {
        return
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    var status string
    var context MessageContext
    if len(user.Emails) > 0 {
        env := config.NewEnvironment()
        context = handler.helper.BuildUserContext(user, params, env, clientToken.Claims["client_id"].(string), handler.guidGenerator, emailTemplate)
        status = handler.helper.SendMailToUser(context, handler.logger, handler.mailClient)
    } else {
        status = "User had no email address"
    }

    response := handler.generateResponse(status, userGUID, context.MessageID)
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}

func (handler NotifyUser) generateResponse(status, userGUID, notificationID string) []byte {
    results := []map[string]string{
        map[string]string{
            "status":          status,
            "recipient":       userGUID,
            "notification_id": notificationID,
        },
    }
    response, err := json.Marshal(results)
    if err != nil {
        panic(err)
    }
    return response
}
