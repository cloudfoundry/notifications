package postal

import (
    "encoding/json"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    EmailFieldName        = "email"
    RecipientsFieldName   = "recipient"
    EmailTextTemplateName = "email_body.text"
    EmailHTMLTemplateName = "email_body.html"
)

type MailRecipeInterface interface {
    Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error)
    Trim([]byte) []byte
}

type EmailRecipe struct {
    mailer         MailerInterface
    templateLoader TemplateLoaderInterface
}

func NewEmailRecipe(mailer MailerInterface, templateLoader TemplateLoaderInterface) EmailRecipe {
    return EmailRecipe{
        mailer:         mailer,
        templateLoader: templateLoader,
    }
}

func (recipe EmailRecipe) Dispatch(clientID string, guid TypedGUID,
    options Options, conn models.ConnectionInterface) ([]Response, error) {

    users := map[string]uaa.User{"no-guid-yet": uaa.User{Emails: []string{options.To}}}
    space := ""
    organization := ""

    subjectTemplate := recipe.determineSubjectTemplate(options.Subject)
    templates, err := recipe.templateLoader.LoadNamedTemplates(subjectTemplate, EmailTextTemplateName, EmailHTMLTemplateName)

    if err != nil {
        return []Response{}, TemplateLoadError("An email template could not be loaded")
    }

    return recipe.mailer.Deliver(conn, templates, users, options, space, organization, clientID), nil
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
