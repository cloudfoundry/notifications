package postal

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
)

const (
    UserContentSuffix         = "user_body"
    OrganizationContentSuffix = "organization_body"
    SpaceContentSuffix        = "space_body"
    SubjectProvidedSuffix     = "subject.provided"
    SubjectMissingSuffix      = "subject.missing"
)

type UserRecipe struct {
    tokenLoader     TokenLoaderInterface
    userLoader      UserLoaderInterface
    templatesLoader TemplatesLoaderInterface
    mailer          MailerInterface
    receiptsRepo    models.ReceiptsRepoInterface
}

func NewUserRecipe(tokenLoader TokenLoaderInterface, userLoader UserLoaderInterface,
    templatesLoader TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) UserRecipe {

    return UserRecipe{
        tokenLoader:     tokenLoader,
        userLoader:      userLoader,
        templatesLoader: templatesLoader,
        mailer:          mailer,
        receiptsRepo:    receiptsRepo,
    }
}

func (recipe UserRecipe) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
    responses := []Response{}

    token, err := recipe.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    users, err := recipe.userLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    subjectSuffix := recipe.subjectSuffix(options.Subject)
    templates, err := recipe.templatesLoader.LoadTemplates(subjectSuffix, UserContentSuffix, clientID, options.KindID)
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

    responses = recipe.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID)

    return responses, nil
}

func (recipe UserRecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, EmailFieldName)
}

func (recipe UserRecipe) subjectSuffix(subject string) string {
    if subject == "" {
        return SubjectMissingSuffix
    }
    return SubjectProvidedSuffix
}
