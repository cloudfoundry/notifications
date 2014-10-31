package postal

import "github.com/cloudfoundry-incubator/notifications/models"

const (
    UserContentSuffix     = "user_body"
    SpaceContentSuffix    = "space_body"
    SubjectProvidedSuffix = "subject.provided"
    SubjectMissingSuffix  = "subject.missing"
)

type UAARecipe struct {
    tokenLoader     TokenLoader
    userLoader      UserLoader
    spaceLoader     SpaceLoader
    templatesLoader TemplatesLoaderInterface
    mailer          MailerInterface
    receiptsRepo    models.ReceiptsRepoInterface
}

func NewUAARecipe(tokenLoader TokenLoader, userLoader UserLoader, spaceLoader SpaceLoader,
    templatesLoader TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) UAARecipe {

    return UAARecipe{
        tokenLoader:     tokenLoader,
        userLoader:      userLoader,
        spaceLoader:     spaceLoader,
        templatesLoader: templatesLoader,
        mailer:          mailer,
        receiptsRepo:    receiptsRepo,
    }
}

func (recipe UAARecipe) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
    responses := []Response{}

    token, err := recipe.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    space, organization, err := recipe.spaceLoader.Load(guid, token)
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

func (recipe UAARecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, EmailFieldName)
}

func (recipe UAARecipe) subjectSuffix(subject string) string {
    if subject == "" {
        return SubjectMissingSuffix
    }
    return SubjectProvidedSuffix
}

func (recipe UAARecipe) contentSuffix(guid TypedGUID) string {
    if guid.BelongsToSpace() {
        return SpaceContentSuffix
    }
    return UserContentSuffix
}
