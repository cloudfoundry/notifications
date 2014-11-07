package postal

import "github.com/cloudfoundry-incubator/notifications/models"

type SpaceStrategy struct {
    tokenLoader       TokenLoaderInterface
    userLoader        UserLoaderInterface
    spaceAndOrgLoader SpaceAndOrgLoaderInterface
    templatesLoader   TemplatesLoaderInterface
    mailer            MailerInterface
    receiptsRepo      models.ReceiptsRepoInterface
}

func NewSpaceStrategy(tokenLoader TokenLoaderInterface, userLoader UserLoaderInterface, spaceAndOrgLoader SpaceAndOrgLoaderInterface,
    templatesLoader TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) SpaceStrategy {

    return SpaceStrategy{
        tokenLoader:       tokenLoader,
        userLoader:        userLoader,
        spaceAndOrgLoader: spaceAndOrgLoader,
        templatesLoader:   templatesLoader,
        mailer:            mailer,
        receiptsRepo:      receiptsRepo,
    }
}

func (strategy SpaceStrategy) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
    responses := []Response{}

    token, err := strategy.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    space, organization, err := strategy.spaceAndOrgLoader.Load(guid, token)
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

func (strategy SpaceStrategy) contentSuffix(guid TypedGUID) string {
    if guid.BelongsToSpace() {
        return models.SpaceBodyTemplateName
    } else if guid.BelongsToOrganization() {
        return models.OrganizationBodyTemplateName
    }
    return models.UserBodyTemplateName
}
