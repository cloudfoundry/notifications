package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
)

type SpaceStrategy struct {
	tokenLoader        utilities.TokenLoaderInterface
	userLoader         utilities.UserLoaderInterface
	spaceLoader        utilities.SpaceLoaderInterface
	organizationLoader utilities.OrganizationLoaderInterface
	findsUserGUIDs     utilities.FindsUserGUIDsInterface
	templatesLoader    utilities.TemplatesLoaderInterface
	mailer             MailerInterface
	receiptsRepo       models.ReceiptsRepoInterface
}

func NewSpaceStrategy(tokenLoader utilities.TokenLoaderInterface, userLoader utilities.UserLoaderInterface, spaceLoader utilities.SpaceLoaderInterface,
	organizationLoader utilities.OrganizationLoaderInterface, findsUserGUIDs utilities.FindsUserGUIDsInterface, templatesLoader utilities.TemplatesLoaderInterface,
	mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) SpaceStrategy {

	return SpaceStrategy{
		tokenLoader:        tokenLoader,
		userLoader:         userLoader,
		spaceLoader:        spaceLoader,
		organizationLoader: organizationLoader,
		findsUserGUIDs:     findsUserGUIDs,
		templatesLoader:    templatesLoader,
		mailer:             mailer,
		receiptsRepo:       receiptsRepo,
	}
}

func (strategy SpaceStrategy) Dispatch(clientID, guid string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	responses := []Response{}

	token, err := strategy.tokenLoader.Load()
	if err != nil {
		return responses, err
	}

	space, err := strategy.spaceLoader.Load(guid, token)
	if err != nil {
		return responses, err
	}

	organization, err := strategy.organizationLoader.Load(space.OrganizationGUID, token)
	if err != nil {
		return responses, err
	}

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToSpace(guid, token)
	if err != nil {
		return responses, err
	}

	users, err := strategy.userLoader.Load(userGUIDs, token)
	if err != nil {
		return responses, err
	}

	subjectSuffix := strategy.subjectSuffix(options.Subject)
	templates, err := strategy.templatesLoader.LoadTemplates(subjectSuffix, models.SpaceBodyTemplateName, clientID, options.KindID)
	if err != nil {
		return responses, postal.TemplateLoadError("An email template could not be loaded")
	}

	err = strategy.receiptsRepo.CreateReceipts(conn, userGUIDs, clientID, options.KindID)
	if err != nil {
		return responses, err
	}

	responses = strategy.mailer.Deliver(conn, templates, users, options, space, organization, clientID, "")

	return responses, nil
}

func (strategy SpaceStrategy) Trim(responses []byte) []byte {
	t := Trimmer{}
	return t.TrimFields(responses, EmailFieldName)
}

func (strategy SpaceStrategy) subjectSuffix(subject string) string {
	if subject == "" {
		return models.SubjectMissingTemplateName
	}
	return models.SubjectProvidedTemplateName
}
