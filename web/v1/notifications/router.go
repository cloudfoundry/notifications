package notifications

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	RequestLogging                   stack.Middleware
	DatabaseAllocator                stack.Middleware
	NotificationsWriteAuthenticator  stack.Middleware
	NotificationsManageAuthenticator stack.Middleware

	Registrar           services.RegistrarInterface
	NotificationsFinder services.NotificationsFinderInterface
	ErrorWriter         handlers.ErrorWriterInterface
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	registerNotificationsHandler := handlers.NewRegisterNotifications(config.Registrar, config.ErrorWriter)
	registerNotificationsStack := stack.NewStack(registerNotificationsHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/registration", registerNotificationsStack).Methods("PUT").Name("PUT /registration")

	registerClientWithNotificationHandler := handlers.NewRegisterClientWithNotifications(config.Registrar, config.ErrorWriter)
	registerClientWithNotificationStack := stack.NewStack(registerClientWithNotificationHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/notifications", registerClientWithNotificationStack).Methods("PUT").Name("PUT /notifications")

	getAllNotificationsHandler := handlers.NewGetAllNotifications(config.NotificationsFinder, config.ErrorWriter)
	getAllNotificationsStack := stack.NewStack(getAllNotificationsHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/notifications", getAllNotificationsStack).Methods("GET").Name("GET /notifications")

	return router
}
