package postal

import (
    "encoding/json"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    EmailFieldName      = "email"
    RecipientsFieldName = "recipient"
    EmptyIDForNonUser   = ""
    EmailSuffix         = "email_body"
)

type MailRecipeInterface interface {
    Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error)
    Trim([]byte) []byte
}

type EmailRecipe struct {
    mailer          MailerInterface
    templatesLoader TemplatesLoaderInterface
}

func NewEmailRecipe(mailer MailerInterface, templatesLoader TemplatesLoaderInterface) EmailRecipe {
    return EmailRecipe{
        mailer:          mailer,
        templatesLoader: templatesLoader,
    }
}

func (recipe EmailRecipe) Dispatch(clientID string, guid TypedGUID,
    options Options, conn models.ConnectionInterface) ([]Response, error) {

    users := map[string]uaa.User{EmptyIDForNonUser: uaa.User{Emails: []string{options.To}}}

    subjectSuffix := recipe.subjectSuffix(options.Subject)

    templates, err := recipe.templatesLoader.LoadTemplates(subjectSuffix, EmailSuffix, clientID, options.KindID)

    if err != nil {
        return []Response{}, TemplateLoadError("An email template could not be loaded")
    }

    return recipe.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID), nil
}

func (recipe EmailRecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, RecipientsFieldName)
}

func (recipe EmailRecipe) determineSubjectTemplate(subject string) string {
    if subject == "" {
        return SubjectMissingTemplateName
    }
    return SubjectProvidedTemplateName
}

type Trimmer struct{}

func (t Trimmer) TrimFields(responses []byte, field string) []byte {
    var results []map[string]string

    err := json.Unmarshal(responses, &results)
    if err != nil {
        panic(err)
    }

    for _, value := range results {
        delete(value, field)
    }

    responses, err = json.Marshal(results)
    if err != nil {
        panic(err)
    }

    return responses
}

func (recipe EmailRecipe) subjectSuffix(subject string) string {
    if subject == "" {
        return SubjectMissingSuffix
    }
    return SubjectProvidedSuffix
}
