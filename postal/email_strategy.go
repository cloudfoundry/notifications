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
)

type StrategyInterface interface {
    Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error)
    Trim([]byte) []byte
}

type EmailStrategy struct {
    mailer          MailerInterface
    templatesLoader TemplatesLoaderInterface
}

func NewEmailStrategy(mailer MailerInterface, templatesLoader TemplatesLoaderInterface) EmailStrategy {
    return EmailStrategy{
        mailer:          mailer,
        templatesLoader: templatesLoader,
    }
}

func (strategy EmailStrategy) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {

    users := map[string]uaa.User{EmptyIDForNonUser: uaa.User{Emails: []string{options.To}}}

    subjectSuffix := strategy.subjectSuffix(options.Subject)

    templates, err := strategy.templatesLoader.LoadTemplates(subjectSuffix, models.EmailBodyTemplateName, clientID, options.KindID)

    if err != nil {
        return []Response{}, TemplateLoadError("An email template could not be loaded")
    }

    return strategy.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID), nil
}

func (strategy EmailStrategy) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, RecipientsFieldName)
}

func (strategy EmailStrategy) determineSubjectTemplate(subject string) string {
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

func (strategy EmailStrategy) subjectSuffix(subject string) string {
    if subject == "" {
        return models.SubjectMissingTemplateName
    }
    return models.SubjectProvidedTemplateName
}
