package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/dgrijalva/jwt-go"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type NotifySpace struct {
    logger          *log.Logger
    cloudController cf.CloudControllerInterface
    uaaClient       uaa.UAAInterface
    mailClient      mail.ClientInterface
    guidGenerator   GUIDGenerationFunc
}

const spaceEmailTemplate = `The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}" component of Cloud Foundry because you are a member of the "{{.Space}}" space in the "{{.Organization}}" organization:

{{.Text}}`

func NewNotifySpace(logger *log.Logger, cloudController cf.CloudControllerInterface,
    uaaClient uaa.UAAInterface, mailClient mail.ClientInterface, guidGenerator GUIDGenerationFunc) NotifySpace {
    return NotifySpace{
        logger:          logger,
        cloudController: cloudController,
        uaaClient:       uaaClient,
        mailClient:      mailClient,
        guidGenerator:   guidGenerator,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    spaceGuid := strings.TrimPrefix(req.URL.Path, "/spaces/")

    params := NewNotifySpaceParams(req.Body)
    if !params.Validate() {
        Error(w, 422, params.Errors)
        return
    }

    token, err := handler.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    handler.uaaClient.SetToken(token.Access)

    ccUsers, err := handler.cloudController.GetUsersBySpaceGuid(spaceGuid, token.Access)
    if err != nil {
        Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    env := config.NewEnvironment()
    space, organization, err := handler.loadSpaceAndOrganization(spaceGuid, token.Access)
    if err != nil {
        Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    for _, ccUser := range ccUsers {
        handler.logger.Println(ccUser.Guid)
        user, ok := loadUser(w, ccUser.Guid, handler.uaaClient)
        if !ok {
            return
        }

        if len(user.Emails) > 0 {
            context := handler.buildContext(user, params, env, space, organization, clientToken.Claims["client_id"].(string))
            status := sendMailToUser(context, handler.logger, handler.mailClient)
            handler.logger.Println(status)
        }
    }
}

func (handler NotifySpace) buildContext(user uaa.User, params NotifySpaceParams, env config.Environment, space, organization, clientID string) MessageContext {
    guid, err := handler.guidGenerator()
    if err != nil {
        panic(err)
    }

    return MessageContext{
        From:              env.Sender,
        To:                user.Emails[0],
        Subject:           params.Subject,
        Text:              params.Text,
        Template:          spaceEmailTemplate,
        KindDescription:   params.KindDescription,
        SourceDescription: params.SourceDescription,
        ClientID:          clientID,
        MessageID:         guid.String(),
        Space:             space,
        Organization:      organization,
    }
}

func (handler NotifySpace) loadSpaceAndOrganization(spaceGuid, token string) (string, string, error) {
    space, err := handler.cloudController.LoadSpace(spaceGuid, token)
    if err != nil {
        return "", "", err
    }

    org, err := handler.cloudController.LoadOrganization(space.OrganizationGuid, token)
    if err != nil {
        return "", "", err
    }

    return space.Name, org.Name, nil
}
