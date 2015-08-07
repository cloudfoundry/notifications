package notify

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
}

type Routes struct {
	RequestLogging                  stack.Middleware
	DatabaseAllocator               stack.Middleware
	NotificationsWriteAuthenticator stack.Middleware
	EmailsWriteAuthenticator        stack.Middleware

	Notify               NotifyInterface
	ErrorWriter          errorWriter
	UserStrategy         services.StrategyInterface
	SpaceStrategy        services.StrategyInterface
	OrganizationStrategy services.StrategyInterface
	EveryoneStrategy     services.StrategyInterface
	UAAScopeStrategy     services.StrategyInterface
	EmailStrategy        services.StrategyInterface
}

func (r Routes) Register(m muxer) {
	requestCounter := middleware.NewRequestCounter(m.GetRouter(), metrics.DefaultLogger)
	m.Handle("POST", "/users/{user_id}", NewUserHandler(r.Notify, r.ErrorWriter, r.UserStrategy), r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/spaces/{space_id}", NewSpaceHandler(r.Notify, r.ErrorWriter, r.SpaceStrategy), r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/organizations/{org_id}", NewOrganizationHandler(r.Notify, r.ErrorWriter, r.OrganizationStrategy), r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/everyone", NewEveryoneHandler(r.Notify, r.ErrorWriter, r.EveryoneStrategy), r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/uaa_scopes/{scope}", NewUAAScopeHandler(r.Notify, r.ErrorWriter, r.UAAScopeStrategy), r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/emails", NewEmailHandler(r.Notify, r.ErrorWriter, r.EmailStrategy), r.RequestLogging, requestCounter, r.EmailsWriteAuthenticator, r.DatabaseAllocator)
}
