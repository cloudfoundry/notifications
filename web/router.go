package web

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/middleware"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/gorilla/mux"
    "github.com/ryanmoran/stack"
)

type MotherInterface interface {
    Registrar() services.Registrar
    EmailRecipe() postal.MailRecipeInterface
    NotificationFinder() services.NotificationFinder
    PreferencesFinder() *services.PreferencesFinder
    PreferenceUpdater() services.PreferenceUpdater
    Logging() stack.Middleware
    ErrorWriter() handlers.ErrorWriter
    Authenticator([]string) middleware.Authenticator
    CORS() middleware.CORS
    handlers.RecipeBuilderInterface
}

type Router struct {
    stacks map[string]stack.Stack
}

func NewRouter(mother MotherInterface) Router {
    registrar := mother.Registrar()
    emailRecipe := mother.EmailRecipe()
    notify := handlers.NewNotify(mother.NotificationFinder(), registrar)
    preferencesFinder := mother.PreferencesFinder()
    preferenceUpdater := mother.PreferenceUpdater()
    logging := mother.Logging()
    errorWriter := mother.ErrorWriter()
    notificationsWriteAuthenticator := mother.Authenticator([]string{"notifications.write"})
    notificationPreferencesReadAuthenticator := mother.Authenticator([]string{"notification_preferences.read"})
    notificationPreferencesWriteAuthenticator := mother.Authenticator([]string{"notification_preferences.write"})
    emailsWriteAuthenticator := mother.Authenticator([]string{"emails.write"})
    cors := mother.CORS()

    return Router{
        stacks: map[string]stack.Stack{
            "GET /info":                 stack.NewStack(handlers.NewGetInfo()).Use(logging),
            "POST /users/{guid}":        stack.NewStack(handlers.NewNotifyUser(notify, errorWriter, mother)).Use(logging, notificationsWriteAuthenticator),
            "POST /spaces/{guid}":       stack.NewStack(handlers.NewNotifySpace(notify, errorWriter, mother)).Use(logging, notificationsWriteAuthenticator),
            "POST /emails":              stack.NewStack(handlers.NewNotifyEmail(notify, errorWriter, emailRecipe)).Use(logging, emailsWriteAuthenticator),
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
