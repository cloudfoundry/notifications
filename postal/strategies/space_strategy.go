package strategies

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type SpaceStrategy struct {
    tokenLoader        postal.TokenLoaderInterface
    userLoader         postal.UserLoaderInterface
    spaceLoader        postal.SpaceLoaderInterface
    organizationLoader postal.OrganizationLoaderInterface
    templatesLoader    postal.TemplatesLoaderInterface
    mailer             MailerInterface
    receiptsRepo       models.ReceiptsRepoInterface
}

func NewSpaceStrategy(tokenLoader postal.TokenLoaderInterface, userLoader postal.UserLoaderInterface, spaceLoader postal.SpaceLoaderInterface,
    organizationLoader postal.OrganizationLoaderInterface, templatesLoader postal.TemplatesLoaderInterface, mailer MailerInterface,
    receiptsRepo models.ReceiptsRepoInterface) SpaceStrategy {

    return SpaceStrategy{
        tokenLoader:        tokenLoader,
        userLoader:         userLoader,
        spaceLoader:        spaceLoader,
        organizationLoader: organizationLoader,
        templatesLoader:    templatesLoader,
        mailer:             mailer,
        receiptsRepo:       receiptsRepo,
    }
}

func (strategy SpaceStrategy) Dispatch(clientID string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
    responses := []Response{}

    token, err := strategy.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    space, err := strategy.spaceLoader.Load(guid.String(), token)
    if err != nil {
        return responses, err
    }

    organization, err := strategy.organizationLoader.Load(space.OrganizationGUID, token)
    if err != nil {
        return responses, err
    }

    users, err := strategy.userLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    subjectSuffix := strategy.subjectSuffix(options.Subject)
    templates, err := strategy.templatesLoader.LoadTemplates(subjectSuffix, models.SpaceBodyTemplateName, clientID, options.KindID)
    if err != nil {
        return responses, postal.TemplateLoadError("An email template could not be loaded")
    }

    var userGUIDs []string
    for key := range users {
        userGUIDs = append(userGUIDs, key)
    }

    err = strategy.receiptsRepo.CreateReceipts(conn, userGUIDs, clientID, options.KindID)
    if err != nil {
        return responses, err
    }

    responses = strategy.mailer.Deliver(conn, templates, users, options, space, organization, clientID)

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
