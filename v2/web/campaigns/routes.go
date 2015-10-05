package campaigns

import (
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging             stack.Middleware
	Authenticator              stack.Middleware
	DatabaseAllocator          stack.Middleware
	CampaignsCollection        collections.CampaignsCollection
	CampaignStatusesCollection collections.CampaignStatusesCollection
	Clock                      clock
}

func (r Routes) Register(m muxer) {
	m.Handle("POST", "/senders/{sender_id}/campaigns", NewCreateHandler(r.CampaignsCollection, r.Clock), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("GET", "/campaigns/{campaign_id}", NewGetHandler(r.CampaignsCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("GET", "/campaigns/{campaign_id}/status", NewStatusHandler(r.CampaignStatusesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
}
