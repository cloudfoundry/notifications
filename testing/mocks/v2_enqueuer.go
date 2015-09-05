package mocks

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
)

type V2Enqueuer struct {
	EnqueueCall struct {
		WasCalled bool
		Receives  struct {
			Connection      queue.ConnectionInterface
			Users           []queue.User
			Options         queue.Options
			Space           cf.CloudControllerSpace
			Org             cf.CloudControllerOrganization
			Client          string
			Scope           string
			VCAPRequestID   string
			RequestReceived time.Time
			UAAHost         string
			CampaignID      string
		}
		Returns struct {
			Responses []queue.Response
		}
	}
}

func NewV2Enqueuer() *V2Enqueuer {
	return &V2Enqueuer{}
}

func (m *V2Enqueuer) Enqueue(conn queue.ConnectionInterface, users []queue.User, options queue.Options,
	space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client, uaaHost, scope, vcapRequestID string, reqReceived time.Time, campaignID string) []queue.Response {

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
	m.EnqueueCall.Receives.CampaignID = campaignID

	m.EnqueueCall.WasCalled = true
	return m.EnqueueCall.Returns.Responses
}
