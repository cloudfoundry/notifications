package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type NotifySpace struct {
    logger          *log.Logger
    cloudController cf.CloudControllerInterface
    uaaClient       uaa.UAAInterface
    mailClient      mail.ClientInterface
}

const spaceEmailTemplate = `The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}" component of Cloud Foundry because you are a member of the "{{.Space}}" space in the "{{.Organization}}" organization:

{{.Text}}`

func NewNotifySpace(logger *log.Logger, cloudController cf.CloudControllerInterface, uaaClient uaa.UAAInterface, mailClient mail.ClientInterface) NotifySpace {
    return NotifySpace{
        logger:          logger,
        cloudController: cloudController,
        uaaClient:       uaaClient,
        mailClient:      mailClient,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    spaceGuid := strings.TrimPrefix(req.URL.Path, "/spaces/")

    params := NewNotifySpaceParams(req.Body)
    if !params.Validate() {
        handler.Error(w, 422, params.Errors)
        return
    }

    token, err := handler.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    handler.uaaClient.SetToken(token.Access)

    ccUsers, err := handler.cloudController.GetUsersBySpaceGuid(spaceGuid, token.Access)
    if err != nil {
        handler.Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    env := config.NewEnvironment()
    space, organization, err := handler.loadSpaceAndOrganization(spaceGuid, token.Access)
    if err != nil {
        handler.Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    for _, ccUser := range ccUsers {
        handler.logger.Println(ccUser.Guid)
        user, ok := handler.loadUser(w, ccUser.Guid)
        if !ok {
            return
        }

        if len(user.Emails) > 0 {
            status := handler.sendMailToUser(user, params, env, space, organization)
            handler.logger.Println(status)
        }
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

func (handler NotifySpace) loadUser(w http.ResponseWriter, guid string) (uaa.User, bool) {
    user, err := handler.uaaClient.UserByID(guid)
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

func (handler NotifySpace) sendMailToUser(user uaa.User, params NotifySpaceParams, env config.Environment, space, organization string) string {
    context := MessageContext{
        From:              env.Sender,
        To:                user.Emails[0],
        Subject:           params.Subject,
        Text:              params.Text,
        Template:          spaceEmailTemplate,
        KindDescription:   params.KindDescription,
        SourceDescription: params.SourceDescription,
        ClientID:          "",
        MessageID:         "",
        Space:             space,
        Organization:      organization,
    }

    status, _, err := SendMail(handler.mailClient, context)
    if err != nil {
        panic(err)
    }

    return status
}

func (handler NotifySpace) Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}
