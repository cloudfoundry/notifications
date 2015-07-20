package notify

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
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

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	userHandler := NewUserHandler(config.Notify, config.ErrorWriter, config.UserStrategy)
	userStack := stack.NewStack(userHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/users/{user_id}", userStack).Methods("POST").Name("POST /users/{user_id}")

	spaceHandler := NewSpaceHandler(config.Notify, config.ErrorWriter, config.SpaceStrategy)
	spaceStack := stack.NewStack(spaceHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/spaces/{space_id}", spaceStack).Methods("POST").Name("POST /spaces/{space_id}")

	orgHandler := NewOrganizationHandler(config.Notify, config.ErrorWriter, config.OrganizationStrategy)
	orgStack := stack.NewStack(orgHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/organizations/{org_id}", orgStack).Methods("POST").Name("POST /organizations/{org_id}")

	everyoneHandler := NewEveryoneHandler(config.Notify, config.ErrorWriter, config.EveryoneStrategy)
	everyoneStack := stack.NewStack(everyoneHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/everyone", everyoneStack).Methods("POST").Name("POST /everyone")

	scopeHandler := NewUAAScopeHandler(config.Notify, config.ErrorWriter, config.UAAScopeStrategy)
	scopeStack := stack.NewStack(scopeHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/uaa_scopes/{scope}", scopeStack).Methods("POST").Name("POST /uaa_scopes/{scope}")

	emailHandler := NewEmailHandler(config.Notify, config.ErrorWriter, config.EmailStrategy)
	emailStack := stack.NewStack(emailHandler).Use(config.RequestLogging, requestCounter, config.EmailsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/emails", emailStack).Methods("POST").Name("POST /emails")

	return router
}
