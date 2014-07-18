package handlers

import (
    "encoding/json"
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
    helper          NotifyHelper
}

const spacePlainTextEmailTemplate = `The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}" component of Cloud Foundry because you are a member of the "{{.Space}}" space in the "{{.Organization}}" organization:

{{.Text}}`

const spaceHTMLEmailTemplate = `<p>The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}" component of Cloud Foundry because you are a member of the "{{.Space}}" space in the "{{.Organization}}" organization:</p>

{{.HTML}}`

func NewNotifySpace(logger *log.Logger, cloudController cf.CloudControllerInterface,
    uaaClient uaa.UAAInterface, mailClient mail.ClientInterface, guidGenerator GUIDGenerationFunc) NotifySpace {
    return NotifySpace{
        logger:          logger,
        cloudController: cloudController,
        uaaClient:       uaaClient,
        mailClient:      mailClient,
        guidGenerator:   guidGenerator,
        helper:          NotifyHelper{},
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    spaceGuid := strings.TrimPrefix(req.URL.Path, "/spaces/")

    params := NewNotifyParams(req.Body)
    if !params.Validate() {
        handler.helper.Error(w, 422, params.Errors)
        return
    }

    token, err := handler.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    handler.uaaClient.SetToken(token.Access)

    ccUsers, err := handler.cloudController.GetUsersBySpaceGuid(spaceGuid, token.Access)
    if err != nil {
        handler.helper.Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    env := config.NewEnvironment()
    space, organization, err := handler.loadSpaceAndOrganization(spaceGuid, token.Access)
    if err != nil {
        handler.helper.Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    responseInformation := make([]map[string]string, len(ccUsers))
    for index, ccUser := range ccUsers {
        handler.logger.Println(ccUser.Guid)
        user, ok := handler.helper.LoadUser(w, ccUser.Guid, handler.uaaClient)
        if !ok {
            return
        }

        if len(user.Emails) > 0 {
            context := handler.helper.BuildSpaceContext(user, params, env, space, organization, clientToken.Claims["client_id"].(string), handler.guidGenerator, spacePlainTextEmailTemplate, spaceHTMLEmailTemplate)
            status := handler.helper.SendMailToUser(context, handler.logger, handler.mailClient)
            handler.logger.Println(status)

            userInfo := make(map[string]string)
            userInfo["status"] = status
            userInfo["recipient"] = ccUser.Guid
            userInfo["notification_id"] = context.MessageID

            responseInformation[index] = userInfo
        }
    }

    response := handler.generateResponse(responseInformation)
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}

func (handler NotifySpace) generateResponse(userInformation []map[string]string) []byte {
    response, err := json.Marshal(userInformation)
    if err != nil {
        panic(err)
    }

    return response
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
