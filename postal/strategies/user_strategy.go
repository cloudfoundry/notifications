package strategies

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type UserStrategy struct {
    tokenLoader     postal.TokenLoaderInterface
    userLoader      postal.UserLoaderInterface
    templatesLoader postal.TemplatesLoaderInterface
    mailer          MailerInterface
    receiptsRepo    models.ReceiptsRepoInterface
}

func NewUserStrategy(tokenLoader postal.TokenLoaderInterface, userLoader postal.UserLoaderInterface,
    templatesLoader postal.TemplatesLoaderInterface, mailer MailerInterface, receiptsRepo models.ReceiptsRepoInterface) UserStrategy {

    return UserStrategy{
        tokenLoader:     tokenLoader,
        userLoader:      userLoader,
        templatesLoader: templatesLoader,
        mailer:          mailer,
        receiptsRepo:    receiptsRepo,
    }
}

func (strategy UserStrategy) Dispatch(clientID string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
    responses := []Response{}

    token, err := strategy.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    users, err := strategy.userLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    subjectSuffix := strategy.subjectSuffix(options.Subject)
    templates, err := strategy.templatesLoader.LoadTemplates(subjectSuffix, models.UserBodyTemplateName, clientID, options.KindID)
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

    responses = strategy.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID)

    return responses, nil
}

func (strategy UserStrategy) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, EmailFieldName)
}

func (strategy UserStrategy) subjectSuffix(subject string) string {
    if subject == "" {
        return models.SubjectMissingTemplateName
    }
    return models.SubjectProvidedTemplateName
}
