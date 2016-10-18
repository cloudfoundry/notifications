package mocks

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type Enqueuer struct {
	EnqueueCall struct {
		WasCalled bool
		Receives  struct {
			Connection      services.ConnectionInterface
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
		Returns struct {
			Responses []services.Response
			Err       error
		}
	}
}

func NewEnqueuer() *Enqueuer {
	return &Enqueuer{}
}

func (m *Enqueuer) Enqueue(
	conn services.ConnectionInterface,
	users []services.User,
	options services.Options,
	space cf.CloudControllerSpace,
	org cf.CloudControllerOrganization,
	client string,
	uaaHost string,
	scope string,
	vcapRequestID string,
	reqReceived time.Time) ([]services.Response, error) {

	m.EnqueueCall.Receives.Connection = conn
	m.EnqueueCall.Receives.Users = users
	m.EnqueueCall.Receives.Options = options
	m.EnqueueCall.Receives.Space = space
	m.EnqueueCall.Receives.Org = org
	m.EnqueueCall.Receives.Client = client
	m.EnqueueCall.Receives.UAAHost = uaaHost
	m.EnqueueCall.Receives.Scope = scope
	m.EnqueueCall.Receives.VCAPRequestID = vcapRequestID
	m.EnqueueCall.Receives.RequestReceived = reqReceived

	m.EnqueueCall.WasCalled = true
	return m.EnqueueCall.Returns.Responses, m.EnqueueCall.Returns.Err
}
