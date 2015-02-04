package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
)

type Mailer struct {
	DeliverArguments map[string]interface{}
	Responses        []strategies.Response
}

func NewMailer() *Mailer {
	return &Mailer{}
}

func (fake *Mailer) Deliver(conn models.ConnectionInterface, users []strategies.User, options postal.Options, space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client, scope string) []strategies.Response {
	fake.DeliverArguments = map[string]interface{}{
		"connection": conn,
		"users":      users,
		"options":    options,
		"space":      space,
		"org":        org,
		"client":     client,
		"scope":      scope,
	}

	return fake.Responses
}
