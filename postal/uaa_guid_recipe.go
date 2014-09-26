package postal

import "github.com/cloudfoundry-incubator/notifications/models"

const (
    SpaceTextTemplateName = "space_body.text"
    SpaceHTMLTemplateName = "space_body.html"
    UserTextTemplateName  = "user_body.text"
    UserHTMLTemplateName  = "user_body.html"
)

type UAARecipe struct {
    tokenLoader    TokenLoader
    userLoader     UserLoader
    spaceLoader    SpaceLoader
    templateLoader TemplateLoader
    mailer         MailerInterface
    receiptsRepo   models.ReceiptsRepoInterface
}

func NewUAARecipe(tokenLoader TokenLoader, userLoader UserLoader, spaceLoader SpaceLoader,
    templateLoader TemplateLoader, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) UAARecipe {
    return UAARecipe{
        tokenLoader:    tokenLoader,
        userLoader:     userLoader,
        spaceLoader:    spaceLoader,
        templateLoader: templateLoader,
        mailer:         mailer,
        receiptsRepo:   receiptsRepo,
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

    subjectTemplate := recipe.subjectTemplate(options.Subject)
    textTemplate := recipe.textTemplate(guid)
    htmlTemplate := recipe.htmlTemplate(guid)

    templates, err := recipe.templateLoader.LoadNamedTemplatesWithClientAndKind(subjectTemplate, textTemplate, htmlTemplate, clientID, options.KindID)
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

func (recipe UAARecipe) subjectTemplate(subject string) string {
    if subject == "" {
        return SubjectMissingTemplateName
    }
    return SubjectProvidedTemplateName
}

func (recipe UAARecipe) textTemplate(guid TypedGUID) string {
    if guid.BelongsToSpace() {
        return SpaceTextTemplateName
    }
    return UserTextTemplateName
}

func (recipe UAARecipe) htmlTemplate(guid TypedGUID) string {
    if guid.BelongsToSpace() {
        return SpaceHTMLTemplateName
    }
    return UserHTMLTemplateName
}
