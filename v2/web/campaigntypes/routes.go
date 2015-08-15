package campaigntypes

import (
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/middleware"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging          stack.Middleware
	Authenticator           middleware.Authenticator
	DatabaseAllocator       middleware.DatabaseAllocator
	CampaignTypesCollection collections.CampaignTypesCollection
}

func (r Routes) Register(m muxer) {
	m.Handle("POST", "/senders/{sender_id}/campaign_types", NewCreateHandler(r.CampaignTypesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("GET", "/senders/{sender_id}/campaign_types", NewListHandler(r.CampaignTypesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("GET", "/senders/{sender_id}/campaign_types/{campaign_type_id:.*}", NewShowHandler(r.CampaignTypesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/senders/{sender_id}/campaign_types/{campaign_type_id}", NewUpdateHandler(r.CampaignTypesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("DELETE", "/senders/{sender_id}/campaign_types/{campaign_type_id}", NewDeleteHandler(r.CampaignTypesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
}
