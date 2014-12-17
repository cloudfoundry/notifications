package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
)

type EveryoneStrategy struct {
	tokenLoader     utilities.TokenLoaderInterface
	allUsers        utilities.AllUsersInterface
	templatesLoader utilities.TemplatesLoaderInterface
	mailer          MailerInterface
	receiptsRepo    models.ReceiptsRepoInterface
}

func NewEveryoneStrategy(tokenLoader utilities.TokenLoaderInterface, allUsers utilities.AllUsersInterface, templatesLoader utilities.TemplatesLoaderInterface, mailer MailerInterface,
	receiptsRepo models.ReceiptsRepoInterface) EveryoneStrategy {
	return EveryoneStrategy{
		tokenLoader:     tokenLoader,
		allUsers:        allUsers,
		templatesLoader: templatesLoader,
		mailer:          mailer,
		receiptsRepo:    receiptsRepo,
	}
}

func (strategy EveryoneStrategy) Dispatch(clientID, guid string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	responses := []Response{}

	_, err := strategy.tokenLoader.Load()
	if err != nil {
		return responses, err
	}

	userEmails, userGUIDs, err := strategy.allUsers.AllUserEmailsAndGUIDs()
	if err != nil {
		return responses, err
	}

	subjectSuffix := strategy.subjectSuffix(options.Subject)
	templates, err := strategy.templatesLoader.LoadTemplates(clientID, options.KindID, models.EveryoneBodyTemplateName, subjectSuffix)
	if err != nil {
		return responses, postal.TemplateLoadError("An email template could not be loaded")
	}

	err = strategy.receiptsRepo.CreateReceipts(conn, userGUIDs, clientID, options.KindID)
	if err != nil {
		return responses, err
	}

	responses = strategy.mailer.Deliver(conn, templates, userEmails, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID, "")

	return responses, nil
}

func (strategy EveryoneStrategy) subjectSuffix(subject string) string {
	if subject == "" {
		return models.SubjectMissingTemplateName
	}
	return models.SubjectProvidedTemplateName
}
