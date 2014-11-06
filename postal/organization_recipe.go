package postal

import "github.com/cloudfoundry-incubator/notifications/models"

type OrganizationRecipe struct {
    tokenLoader       TokenLoaderInterface
    userLoader        UserLoaderInterface
    spaceAndOrgLoader SpaceAndOrgLoaderInterface
    templatesLoader   TemplatesLoaderInterface
    mailer            MailerInterface
    receiptsRepo      models.ReceiptsRepoInterface
}

func NewOrganizationRecipe(tokenLoader TokenLoaderInterface, userLoader UserLoaderInterface, spaceAndOrgLoader SpaceAndOrgLoaderInterface,
    templatesLoader TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) OrganizationRecipe {

    return OrganizationRecipe{
        tokenLoader:       tokenLoader,
        userLoader:        userLoader,
        spaceAndOrgLoader: spaceAndOrgLoader,
        templatesLoader:   templatesLoader,
        mailer:            mailer,
        receiptsRepo:      receiptsRepo,
    }
}

func (recipe OrganizationRecipe) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
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

func (recipe OrganizationRecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, EmailFieldName)
}

func (recipe OrganizationRecipe) subjectSuffix(subject string) string {
    if subject == "" {
        return SubjectMissingSuffix
    }
    return SubjectProvidedSuffix
}

func (recipe OrganizationRecipe) contentSuffix(guid TypedGUID) string {
    if guid.BelongsToSpace() {
        return SpaceContentSuffix
    } else if guid.BelongsToOrganization() {
        return OrganizationContentSuffix
    }
    return UserContentSuffix
}
