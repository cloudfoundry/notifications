package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
)

type UserStrategy struct {
	tokenLoader     utilities.TokenLoaderInterface
	userLoader      utilities.UserLoaderInterface
	templatesLoader utilities.TemplatesLoaderInterface
	mailer          MailerInterface
	receiptsRepo    models.ReceiptsRepoInterface
}

func NewUserStrategy(tokenLoader utilities.TokenLoaderInterface, userLoader utilities.UserLoaderInterface,
	templatesLoader utilities.TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) UserStrategy {

	return UserStrategy{
		tokenLoader:     tokenLoader,
		userLoader:      userLoader,
		templatesLoader: templatesLoader,
		mailer:          mailer,
		receiptsRepo:    receiptsRepo,
	}
}

func (strategy UserStrategy) Dispatch(clientID, guid string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	responses := []Response{}

	token, err := strategy.tokenLoader.Load()
	if err != nil {
		return responses, err
	}

	userGUIDs := []string{guid}
	users, err := strategy.userLoader.Load(userGUIDs, token)
	if err != nil {
		return responses, err
	}

	templates, err := strategy.templatesLoader.LoadTemplates(clientID, options.KindID)
	if err != nil {
		return responses, postal.TemplateLoadError("An email template could not be loaded")
	}

	err = strategy.receiptsRepo.CreateReceipts(conn, userGUIDs, clientID, options.KindID)
	if err != nil {
		return responses, err
	}

	responses = strategy.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID, "")

	return responses, nil
}
