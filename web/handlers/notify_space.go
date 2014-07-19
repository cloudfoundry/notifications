package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/mail"
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
        helper:          NewNotifyHelper(cloudController, logger, uaaClient, guidGenerator, mailClient),
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {

    determineGUID := func(path string) string {
        return strings.TrimPrefix(path, "/spaces/")
    }

    loadUsers := func(spaceGuid, accessToken string) ([]cf.CloudControllerUser, error) {
        ccUsers, err := handler.cloudController.GetUsersBySpaceGuid(spaceGuid, accessToken)
        return ccUsers, err
    }
    loadSpaceAndOrganization := true
    handler.helper.SendMail(w, req, determineGUID, loadUsers, loadSpaceAndOrganization,
        spacePlainTextEmailTemplate, spaceHTMLEmailTemplate)
}
