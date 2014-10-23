package web

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
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
    TemplateFinder() services.TemplateFinder
    Database() models.DatabaseInterface
    Logging() stack.Middleware
    ErrorWriter() handlers.ErrorWriter
    Authenticator(...string) middleware.Authenticator
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
    templateFinder := mother.TemplateFinder()
    logging := mother.Logging()
    errorWriter := mother.ErrorWriter()
    notificationsWriteAuthenticator := mother.Authenticator("notifications.write")
    notificationPreferencesReadAuthenticator := mother.Authenticator("notification_preferences.read")
    notificationPreferencesWriteAuthenticator := mother.Authenticator("notification_preferences.write")
    notificationPreferencesAdminAuthenticator := mother.Authenticator("notification_preferences.admin")
    emailsWriteAuthenticator := mother.Authenticator("emails.write")
    database := mother.Database()

    cors := mother.CORS()

    return Router{
        stacks: map[string]stack.Stack{
            "GET /info":                        stack.NewStack(handlers.NewGetInfo()).Use(logging),
            "POST /users/{guid}":               stack.NewStack(handlers.NewNotifyUser(notify, errorWriter, mother, database)).Use(logging, notificationsWriteAuthenticator),
            "POST /spaces/{guid}":              stack.NewStack(handlers.NewNotifySpace(notify, errorWriter, mother, database)).Use(logging, notificationsWriteAuthenticator),
            "POST /emails":                     stack.NewStack(handlers.NewNotifyEmail(notify, errorWriter, emailRecipe, database)).Use(logging, emailsWriteAuthenticator),
            "PUT /registration":                stack.NewStack(handlers.NewRegisterNotifications(registrar, errorWriter, database)).Use(logging, notificationsWriteAuthenticator),
            "OPTIONS /user_preferences":        stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, cors),
            "OPTIONS /user_preferences/{guid}": stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, cors),
            "GET /user_preferences":            stack.NewStack(handlers.NewGetPreferences(preferencesFinder, errorWriter)).Use(logging, cors, notificationPreferencesReadAuthenticator),
            "GET /user_preferences/{guid}":     stack.NewStack(handlers.NewGetPreferencesForUser(preferencesFinder, errorWriter)).Use(logging, cors, notificationPreferencesAdminAuthenticator),
            "PATCH /user_preferences":          stack.NewStack(handlers.NewUpdatePreferences(preferenceUpdater, errorWriter, database)).Use(logging, cors, notificationPreferencesWriteAuthenticator),
            "PATCH /user_preferences/{guid}":   stack.NewStack(handlers.NewUpdateSpecificUserPreferences(preferenceUpdater, errorWriter, database)).Use(logging, cors, notificationPreferencesAdminAuthenticator),
            "GET /templates/{templateName}":    stack.NewStack(handlers.NewGetTemplates(templateFinder)).Use(logging),
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
