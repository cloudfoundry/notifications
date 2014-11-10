package strategies

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/postal/utilities"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    EmailFieldName      = "email"
    RecipientsFieldName = "recipient"
    EmptyIDForNonUser   = ""
)

type EmailStrategy struct {
    mailer          MailerInterface
    templatesLoader utilities.TemplatesLoaderInterface
}

func NewEmailStrategy(mailer MailerInterface, templatesLoader utilities.TemplatesLoaderInterface) EmailStrategy {
    return EmailStrategy{
        mailer:          mailer,
        templatesLoader: templatesLoader,
    }
}

func (strategy EmailStrategy) Dispatch(clientID, guid string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
    users := map[string]uaa.User{
        EmptyIDForNonUser: uaa.User{
            Emails: []string{options.To},
        },
    }

    templates, err := strategy.templatesLoader.LoadTemplates(strategy.subjectSuffix(options.Subject), models.EmailBodyTemplateName, clientID, options.KindID)
    if err != nil {
        return []Response{}, postal.TemplateLoadError("An email template could not be loaded")
    }

    return strategy.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID), nil
}

func (strategy EmailStrategy) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, RecipientsFieldName)
}

func (strategy EmailStrategy) determineSubjectTemplate(subject string) string {
    if subject == "" {
        return models.SubjectMissingTemplateName
    }
    return models.SubjectProvidedTemplateName
}

func (strategy EmailStrategy) subjectSuffix(subject string) string {
    if subject == "" {
        return models.SubjectMissingTemplateName
    }
    return models.SubjectProvidedTemplateName
}
