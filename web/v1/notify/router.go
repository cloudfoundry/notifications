package notify

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	RequestLogging                  middleware.RequestLogging
	DatabaseAllocator               middleware.DatabaseAllocator
	NotificationsWriteAuthenticator middleware.Authenticator
	EmailsWriteAuthenticator        middleware.Authenticator

	Notify               handlers.NotifyInterface
	ErrorWriter          handlers.ErrorWriterInterface
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

	notifyUserHandler := handlers.NewNotifyUser(config.Notify, config.ErrorWriter, config.UserStrategy)
	notifyUserStack := stack.NewStack(notifyUserHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/users/{user_id}", notifyUserStack).Methods("POST").Name("POST /users/{user_id}")

	notifySpaceHandler := handlers.NewNotifySpace(config.Notify, config.ErrorWriter, config.SpaceStrategy)
	notifySpaceStack := stack.NewStack(notifySpaceHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/spaces/{space_id}", notifySpaceStack).Methods("POST").Name("POST /spaces/{space_id}")

	notifyOrganizationHandler := handlers.NewNotifyOrganization(config.Notify, config.ErrorWriter, config.OrganizationStrategy)
	notifyOrganizationStack := stack.NewStack(notifyOrganizationHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/organizations/{org_id}", notifyOrganizationStack).Methods("POST").Name("POST /organizations/{org_id}")

	notifyEveryoneHandler := handlers.NewNotifyEveryone(config.Notify, config.ErrorWriter, config.EveryoneStrategy)
	notifyEveryoneStack := stack.NewStack(notifyEveryoneHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/everyone", notifyEveryoneStack).Methods("POST").Name("POST /everyone")

	notifyUAAScopeHandler := handlers.NewNotifyUAAScope(config.Notify, config.ErrorWriter, config.UAAScopeStrategy)
	notifyUAAScopeStack := stack.NewStack(notifyUAAScopeHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/uaa_scopes/{scope}", notifyUAAScopeStack).Methods("POST").Name("POST /uaa_scopes/{scope}")

	notifyEmailHandler := handlers.NewNotifyEmail(config.Notify, config.ErrorWriter, config.EmailStrategy)
	notifyEmailStack := stack.NewStack(notifyEmailHandler).Use(config.RequestLogging, requestCounter, config.EmailsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/emails", notifyEmailStack).Methods("POST").Name("POST /emails")

	return router
}
