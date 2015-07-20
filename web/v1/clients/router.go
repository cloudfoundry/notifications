package clients

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	RequestLogging                   stack.Middleware
	NotificationsManageAuthenticator stack.Middleware
	DatabaseAllocator                stack.Middleware

	ErrorWriter      errorWriter
	TemplateAssigner services.TemplateAssignerInterface
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	assignTemplateHandler := NewAssignTemplateHandler(config.TemplateAssigner, config.ErrorWriter)
	assignTemplateStack := stack.NewStack(assignTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/clients/{client_id}/template", assignTemplateStack).Methods("PUT").Name("PUT /clients/{client_id}/template")

	return router
}
