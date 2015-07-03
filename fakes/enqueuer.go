package fakes

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
)

type Enqueuer struct {
	EnqueueCall struct {
		Args struct {
			Connection      models.ConnectionInterface
			Users           []services.User
			Options         services.Options
			Space           cf.CloudControllerSpace
			Org             cf.CloudControllerOrganization
			Client          string
			Scope           string
			VCAPRequestID   string
			RequestReceived time.Time
			UAAHost         string
		}
		Responses []services.Response
	}
}

func NewEnqueuer() *Enqueuer {
	return &Enqueuer{}
}

func (m *Enqueuer) Enqueue(conn models.ConnectionInterface, users []services.User, options services.Options,
	space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client, uaaHost, scope, vcapRequestID string, reqReceived time.Time) []services.Response {

	m.EnqueueCall.Args.Connection = conn
	m.EnqueueCall.Args.Users = users
	m.EnqueueCall.Args.Options = options
	m.EnqueueCall.Args.Space = space
	m.EnqueueCall.Args.Org = org
	m.EnqueueCall.Args.Client = client
	m.EnqueueCall.Args.UAAHost = uaaHost
	m.EnqueueCall.Args.Scope = scope
	m.EnqueueCall.Args.VCAPRequestID = vcapRequestID
	m.EnqueueCall.Args.RequestReceived = reqReceived

	return m.EnqueueCall.Responses
}
