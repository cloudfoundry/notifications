package clients

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Routes struct {
	RequestLogging                   stack.Middleware
	NotificationsManageAuthenticator stack.Middleware
	DatabaseAllocator                stack.Middleware

	ErrorWriter      errorWriter
	TemplateAssigner services.TemplateAssignerInterface
}

func (r Routes) Register(router *mux.Router) {
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	assignTemplateHandler := NewAssignTemplateHandler(r.TemplateAssigner, r.ErrorWriter)
	assignTemplateStack := stack.NewStack(assignTemplateHandler).Use(r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	router.Handle("/clients/{client_id}/template", assignTemplateStack).Methods("PUT").Name("PUT /clients/{client_id}/template")
}
