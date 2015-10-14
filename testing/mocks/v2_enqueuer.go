package mocks

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
)

type V2Enqueuer struct {
	EnqueueCallsCount int
	EnqueueCalls      []V2EnqueuerEnqueueCall
}

type V2EnqueuerEnqueueCall struct {
	Receives V2EnqueuerEnqueueCallReceives
}

type V2EnqueuerEnqueueCallReceives struct {
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

func NewV2Enqueuer() *V2Enqueuer {
	return &V2Enqueuer{}
}

func (m *V2Enqueuer) Enqueue(conn queue.ConnectionInterface, users []queue.User, options queue.Options,
	space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client, uaaHost, scope, vcapRequestID string, reqReceived time.Time, campaignID string) {

	if len(m.EnqueueCalls) <= m.EnqueueCallsCount {
		m.EnqueueCalls = append(m.EnqueueCalls, V2EnqueuerEnqueueCall{})
	}

	call := m.EnqueueCalls[m.EnqueueCallsCount]
	call.Receives.Connection = conn
	call.Receives.Users = users
	call.Receives.Options = options
	call.Receives.Space = space
	call.Receives.Org = org
	call.Receives.Client = client
	call.Receives.UAAHost = uaaHost
	call.Receives.Scope = scope
	call.Receives.VCAPRequestID = vcapRequestID
	call.Receives.RequestReceived = reqReceived
	call.Receives.CampaignID = campaignID
	m.EnqueueCalls[m.EnqueueCallsCount] = call

	m.EnqueueCallsCount++
}
