package postal

import "github.com/cloudfoundry-incubator/notifications/models"

type SpaceRecipe struct {
    tokenLoader       TokenLoaderInterface
    userLoader        UserLoaderInterface
    spaceAndOrgLoader SpaceAndOrgLoaderInterface
    templatesLoader   TemplatesLoaderInterface
    mailer            MailerInterface
    receiptsRepo      models.ReceiptsRepoInterface
}

func NewSpaceRecipe(tokenLoader TokenLoaderInterface, userLoader UserLoaderInterface, spaceAndOrgLoader SpaceAndOrgLoaderInterface,
    templatesLoader TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) SpaceRecipe {

    return SpaceRecipe{
        tokenLoader:       tokenLoader,
        userLoader:        userLoader,
        spaceAndOrgLoader: spaceAndOrgLoader,
        templatesLoader:   templatesLoader,
        mailer:            mailer,
        receiptsRepo:      receiptsRepo,
    }
}

func (recipe SpaceRecipe) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
    responses := []Response{}

    token, err := recipe.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    space, organization, err := recipe.spaceAndOrgLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    users, err := recipe.userLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    subjectSuffix := recipe.subjectSuffix(options.Subject)
    contentSuffix := recipe.contentSuffix(guid)

    templates, err := recipe.templatesLoader.LoadTemplates(subjectSuffix, contentSuffix, clientID, options.KindID)
    if err != nil {
        return responses, TemplateLoadError("An email template could not be loaded")
    }

    var userGUIDs []string
    for key := range users {
        userGUIDs = append(userGUIDs, key)
    }

    err = recipe.receiptsRepo.CreateReceipts(conn, userGUIDs, clientID, options.KindID)
    if err != nil {
        return responses, err
    }

    responses = recipe.mailer.Deliver(conn, templates, users, options, space, organization, clientID)

    return responses, nil
}

func (recipe SpaceRecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, EmailFieldName)
}

func (recipe SpaceRecipe) subjectSuffix(subject string) string {
    if subject == "" {
        return models.SubjectMissingTemplateName
    }
    return models.SubjectProvidedTemplateName
}

func (recipe SpaceRecipe) contentSuffix(guid TypedGUID) string {
    if guid.BelongsToSpace() {
        return models.SpaceBodyTemplateName
    } else if guid.BelongsToOrganization() {
        return models.OrganizationBodyTemplateName
    }
    return models.UserBodyTemplateName
}
