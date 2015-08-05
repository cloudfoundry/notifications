package campaigntypes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	RequestLogging          stack.Middleware
	Authenticator           middleware.Authenticator
	DatabaseAllocator       middleware.DatabaseAllocator
	CampaignTypesCollection collections.CampaignTypesCollection
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()

	createStack := stack.NewStack(NewCreateHandler(config.CampaignTypesCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)
	listStack := stack.NewStack(NewListHandler(config.CampaignTypesCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)
	showStack := stack.NewStack(NewShowHandler(config.CampaignTypesCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)
	updateStack := stack.NewStack(NewUpdateHandler(config.CampaignTypesCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)

	router.Handle("/senders/{sender_id}/campaign_types", createStack).Methods("POST").Name("POST /senders/{sender_id}/campaign_types")
	router.Handle("/senders/{sender_id}/campaign_types", listStack).Methods("GET").Name("GET /senders/{sender_id}/campaign_types")
	router.Handle("/senders/{sender_id}/campaign_types/{campaign_type_id}", showStack).Methods("GET").Name("GET /senders/{sender_id}/campaign_types/{campaign_type_id}")
	router.Handle("/senders/{sender_id}/campaign_types/{campaign_type_id}", updateStack).Methods("PUT").Name("PUT /senders/{sender_id}/campaign_types/{campaign_type_id}")

	return router
}
