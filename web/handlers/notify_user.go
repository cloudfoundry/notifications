package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/mail"
    uuid "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const plainTextEmailTemplate = `The following "{{.KindDescription}}" notification was sent to you directly by the "{{.SourceDescription}}" component of Cloud Foundry:

{{.Text}}`

const htmlEmailTemplate = `<p>The following "{{.KindDescription}}" notification was sent to you directly by the "{{.SourceDescription}}" component of Cloud Foundry:</p>

{{.HTML}}`

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
        helper:        NewNotifyHelper(cf.CloudController{}, logger, uaaClient, guidGenerator, mailClient),
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    determineGUID := func(path string) string {
        return strings.TrimPrefix(path, "/users/")
    }

    loadUsers := func(userGuid, accessToken string) ([]cf.CloudControllerUser, error) {
        ccUsers := []cf.CloudControllerUser{}

        user := cf.CloudControllerUser{
            Guid: userGuid,
        }

        ccUsers = append(ccUsers, user)
        return ccUsers, nil
    }

    loadSpaceAndOrganization := false
    handler.helper.SendMail(w, req, determineGUID, loadUsers, loadSpaceAndOrganization,
        plainTextEmailTemplate, htmlEmailTemplate)
}
