package fakes

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
)

type Enqueuer struct {
	EnqueueCall struct {
		Args struct {
			Connection      models.ConnectionInterface
			Users           []strategies.User
			Options         postal.Options
			Space           cf.CloudControllerSpace
			Org             cf.CloudControllerOrganization
			Client          string
			Scope           string
			VCAPRequestID   string
			RequestReceived time.Time
		}
		Responses []strategies.Response
	}
}

func NewEnqueuer() *Enqueuer {
	return &Enqueuer{}
}

func (m *Enqueuer) Enqueue(conn models.ConnectionInterface, users []strategies.User, options postal.Options,
	space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client, scope, vcapRequestID string, reqReceived time.Time) []strategies.Response {

	m.EnqueueCall.Args.Connection = conn
	m.EnqueueCall.Args.Users = users
	m.EnqueueCall.Args.Options = options
	m.EnqueueCall.Args.Space = space
	m.EnqueueCall.Args.Org = org
	m.EnqueueCall.Args.Client = client
	m.EnqueueCall.Args.Scope = scope
	m.EnqueueCall.Args.VCAPRequestID = vcapRequestID
	m.EnqueueCall.Args.RequestReceived = reqReceived

	return m.EnqueueCall.Responses
}
