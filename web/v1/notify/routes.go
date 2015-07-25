package notify

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

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

func (r Routes) Register(router *mux.Router) {
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	userHandler := NewUserHandler(r.Notify, r.ErrorWriter, r.UserStrategy)
	userStack := stack.NewStack(userHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/users/{user_id}", userStack).Methods("POST").Name("POST /users/{user_id}")

	spaceHandler := NewSpaceHandler(r.Notify, r.ErrorWriter, r.SpaceStrategy)
	spaceStack := stack.NewStack(spaceHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/spaces/{space_id}", spaceStack).Methods("POST").Name("POST /spaces/{space_id}")

	orgHandler := NewOrganizationHandler(r.Notify, r.ErrorWriter, r.OrganizationStrategy)
	orgStack := stack.NewStack(orgHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/organizations/{org_id}", orgStack).Methods("POST").Name("POST /organizations/{org_id}")

	everyoneHandler := NewEveryoneHandler(r.Notify, r.ErrorWriter, r.EveryoneStrategy)
	everyoneStack := stack.NewStack(everyoneHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/everyone", everyoneStack).Methods("POST").Name("POST /everyone")

	scopeHandler := NewUAAScopeHandler(r.Notify, r.ErrorWriter, r.UAAScopeStrategy)
	scopeStack := stack.NewStack(scopeHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/uaa_scopes/{scope}", scopeStack).Methods("POST").Name("POST /uaa_scopes/{scope}")

	emailHandler := NewEmailHandler(r.Notify, r.ErrorWriter, r.EmailStrategy)
	emailStack := stack.NewStack(emailHandler).Use(r.RequestLogging, requestCounter, r.EmailsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/emails", emailStack).Methods("POST").Name("POST /emails")
}
