package web

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/middleware"
    "github.com/gorilla/mux"
    "github.com/ryanmoran/stack"
)

type Router struct {
    stacks map[string]stack.Stack
}

func NewRouter(mother *Mother) Router {
    registrar := mother.Registrar()
    notify := handlers.NewNotify(mother.Courier(), mother.NotificationFinder(), registrar)
    preferencesFinder := mother.PreferencesFinder()
    preferenceUpdater := mother.PreferenceUpdater()
    logging := mother.Logging()
    errorWriter := mother.ErrorWriter()
    notificationsWriteAuthenticator := mother.Authenticator([]string{"notifications.write"})
    notificationPreferencesReadAuthenticator := mother.Authenticator([]string{"notification_preferences.read"})
    notificationPreferencesWriteAuthenticator := mother.Authenticator([]string{"notification_preferences.write"})
    cors := middleware.NewCORS()

    return Router{
        stacks: map[string]stack.Stack{
            "GET /info":                 stack.NewStack(handlers.NewGetInfo()).Use(logging),
            "POST /users/{guid}":        stack.NewStack(handlers.NewNotifyUser(notify, errorWriter)).Use(logging, notificationsWriteAuthenticator),
            "POST /spaces/{guid}":       stack.NewStack(handlers.NewNotifySpace(notify, errorWriter)).Use(logging, notificationsWriteAuthenticator),
            "PUT /registration":         stack.NewStack(handlers.NewRegisterNotifications(registrar, errorWriter)).Use(logging, notificationsWriteAuthenticator),
            "OPTIONS /user_preferences": stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, cors),
            "GET /user_preferences":     stack.NewStack(handlers.NewGetPreferences(preferencesFinder, errorWriter)).Use(logging, cors, notificationPreferencesReadAuthenticator),
            "PATCH /user_preferences":   stack.NewStack(handlers.NewUpdatePreferences(preferenceUpdater, errorWriter)).Use(logging, cors, notificationPreferencesWriteAuthenticator),
        },
    }
}

func (router Router) Routes() *mux.Router {
    r := mux.NewRouter()
    for methodPath, stack := range router.stacks {
        var name = methodPath
        parts := strings.SplitN(methodPath, " ", 2)
        r.Handle(parts[1], stack).Methods(parts[0]).Name(name)
    }
    return r
}
