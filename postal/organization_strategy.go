package postal

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
)

type OrganizationStrategy struct {
    tokenLoader        TokenLoaderInterface
    userLoader         UserLoaderInterface
    organizationLoader OrganizationLoaderInterface
    templatesLoader    TemplatesLoaderInterface
    mailer             MailerInterface
    receiptsRepo       models.ReceiptsRepoInterface
}

func NewOrganizationStrategy(tokenLoader TokenLoaderInterface, userLoader UserLoaderInterface, organizationLoader OrganizationLoaderInterface,
    templatesLoader TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) OrganizationStrategy {

    return OrganizationStrategy{
        tokenLoader:        tokenLoader,
        userLoader:         userLoader,
        organizationLoader: organizationLoader,
        templatesLoader:    templatesLoader,
        mailer:             mailer,
        receiptsRepo:       receiptsRepo,
    }
}

func (strategy OrganizationStrategy) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
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
    contentSuffix := strategy.contentSuffix(guid)

    templates, err := strategy.templatesLoader.LoadTemplates(subjectSuffix, contentSuffix, clientID, options.KindID)
    if err != nil {
        return responses, TemplateLoadError("An email template could not be loaded")
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

func (strategy OrganizationStrategy) contentSuffix(guid TypedGUID) string {
    if guid.BelongsToSpace() {
        return models.SpaceBodyTemplateName
    } else if guid.BelongsToOrganization() {
        return models.OrganizationBodyTemplateName
    }
    return models.UserBodyTemplateName
}
