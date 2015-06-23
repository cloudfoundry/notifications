package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewNotifyRouter(notify handlers.NotifyInterface,
	errorWriter handlers.ErrorWriterInterface,
	userStrategy strategies.StrategyInterface,
	logging middleware.RequestLogging,
	notificationsWriteAuthenticator middleware.Authenticator,
	databaseAllocator middleware.DatabaseAllocator,
	spaceStrategy strategies.StrategyInterface,
	organizationStrategy strategies.StrategyInterface,
	everyoneStrategy strategies.StrategyInterface,
	uaaScopeStrategy strategies.StrategyInterface,
	emailStrategy strategies.StrategyInterface,
	emailsWriteAuthenticator middleware.Authenticator) *mux.Router {

	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/users/{user_id}", stack.NewStack(handlers.NewNotifyUser(notify, errorWriter, userStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /users/{user_id}")
	router.Handle("/spaces/{space_id}", stack.NewStack(handlers.NewNotifySpace(notify, errorWriter, spaceStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /spaces/{space_id}")
	router.Handle("/organizations/{org_id}", stack.NewStack(handlers.NewNotifyOrganization(notify, errorWriter, organizationStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /organizations/{org_id}")
	router.Handle("/everyone", stack.NewStack(handlers.NewNotifyEveryone(notify, errorWriter, everyoneStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /everyone")
	router.Handle("/uaa_scopes/{scope}", stack.NewStack(handlers.NewNotifyUAAScope(notify, errorWriter, uaaScopeStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /uaa_scopes/{scope}")
	router.Handle("/emails", stack.NewStack(handlers.NewNotifyEmail(notify, errorWriter, emailStrategy)).Use(logging, requestCounter, emailsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /emails")

	return router
}
