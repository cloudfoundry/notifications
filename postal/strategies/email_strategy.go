package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const EmailEndorsement = "This message was sent directly to your email address."

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
		options.To: uaa.User{
			Emails: []string{options.To},
		},
	}

	templates, err := strategy.templatesLoader.LoadTemplates(clientID, options.KindID)
	if err != nil {
		return []Response{}, postal.TemplateLoadError("An email template could not be loaded")
	}

	options.Endorsement = EmailEndorsement

	return strategy.mailer.Deliver(conn, templates, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID, ""), nil
}
