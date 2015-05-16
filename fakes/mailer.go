package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
)

type Mailer struct {
	DeliverCall struct {
		Args struct {
			Connection    models.ConnectionInterface
			Users         []strategies.User
			Options       postal.Options
			Space         cf.CloudControllerSpace
			Org           cf.CloudControllerOrganization
			Client        string
			Scope         string
			VCAPRequestID string
		}
		Responses []strategies.Response
	}
}

func NewMailer() *Mailer {
	return &Mailer{}
}

func (m *Mailer) Deliver(conn models.ConnectionInterface, users []strategies.User, options postal.Options,
	space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client, scope, vcapRequestID string) []strategies.Response {

	m.DeliverCall.Args.Connection = conn
	m.DeliverCall.Args.Users = users
	m.DeliverCall.Args.Options = options
	m.DeliverCall.Args.Space = space
	m.DeliverCall.Args.Org = org
	m.DeliverCall.Args.Client = client
	m.DeliverCall.Args.Scope = scope
	m.DeliverCall.Args.VCAPRequestID = vcapRequestID

	return m.DeliverCall.Responses
}
