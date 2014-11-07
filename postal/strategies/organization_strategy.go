package strategies

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type OrganizationStrategy struct {
    tokenLoader        postal.TokenLoaderInterface
    userLoader         postal.UserLoaderInterface
    organizationLoader postal.OrganizationLoaderInterface
    templatesLoader    postal.TemplatesLoaderInterface
    mailer             MailerInterface
    receiptsRepo       models.ReceiptsRepoInterface
}

func NewOrganizationStrategy(tokenLoader postal.TokenLoaderInterface, userLoader postal.UserLoaderInterface, organizationLoader postal.OrganizationLoaderInterface,
    templatesLoader postal.TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) OrganizationStrategy {

    return OrganizationStrategy{
        tokenLoader:        tokenLoader,
        userLoader:         userLoader,
        organizationLoader: organizationLoader,
        templatesLoader:    templatesLoader,
        mailer:             mailer,
        receiptsRepo:       receiptsRepo,
    }
}

func (strategy OrganizationStrategy) Dispatch(clientID string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
    responses := []Response{}

    token, err := strategy.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    organization, err := strategy.organizationLoader.Load(guid.String(), token)
    if err != nil {
        return responses, err
    }

    users, err := strategy.userLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    subjectSuffix := strategy.subjectSuffix(options.Subject)

    templates, err := strategy.templatesLoader.LoadTemplates(subjectSuffix, models.OrganizationBodyTemplateName, clientID, options.KindID)
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

    responses = strategy.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, organization, clientID)

    return responses, nil
}

func (strategy OrganizationStrategy) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, EmailFieldName)
}

func (strategy OrganizationStrategy) subjectSuffix(subject string) string {
    if subject == "" {
        return models.SubjectMissingTemplateName
    }
    return models.SubjectProvidedTemplateName
}
