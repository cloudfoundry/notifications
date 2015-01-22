package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type RegisterClientWithNotifications struct {
	registrar   services.RegistrarInterface
	errorWriter ErrorWriterInterface
	database    models.DatabaseInterface
}

func NewRegisterClientWithNotifications(registrar services.RegistrarInterface, errorWriter ErrorWriterInterface, database models.DatabaseInterface) RegisterClientWithNotifications {
	return RegisterClientWithNotifications{
		registrar:   registrar,
		errorWriter: errorWriter,
		database:    database,
	}
}

func (handler RegisterClientWithNotifications) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.registration",
	}).Log()

	handler.Execute(w, req, handler.database.Connection(), context)
}

func (handler RegisterClientWithNotifications) Execute(w http.ResponseWriter, req *http.Request,
	connection models.ConnectionInterface, context stack.Context) {

	parameters, err := params.NewClientRegistration(req.Body)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	err = parameters.Validate()
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	generatedKinds := []models.Kind{}
	for _, notification := range parameters.Notifications {
		generatedKinds = append(generatedKinds, models.Kind{
			ID:          notification.ID,
			Description: notification.Description,
			Critical:    notification.Critical,
			TemplateID:  models.DoNotSetTemplateID,
		})
	}

	token := context.Get("token").(*jwt.Token)
	clientID := token.Claims["client_id"].(string)

	client := models.Client{
		ID:          clientID,
		Description: parameters.SourceName,
		TemplateID:  models.DoNotSetTemplateID,
	}

	kinds, err := handler.ValidateCriticalScopes(token.Claims["scope"], generatedKinds, client)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	transaction := connection.Transaction()
	transaction.Begin()

	err = handler.registrar.Register(transaction, client, kinds)
	if err != nil {
		transaction.Rollback()
		handler.errorWriter.Write(w, err)
		return
	}

	if len(parameters.Notifications) > 0 {
		err = handler.registrar.Prune(transaction, client, kinds)
		if err != nil {
			transaction.Rollback()
			handler.errorWriter.Write(w, err)
			return
		}
	}

	err = transaction.Commit()
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler RegisterClientWithNotifications) ValidateCriticalScopes(scopes interface{}, kinds []models.Kind, client models.Client) ([]models.Kind, error) {
	hasCriticalWrite := false
	for _, scope := range scopes.([]interface{}) {
		if scope.(string) == "critical_notifications.write" {
			hasCriticalWrite = true
		}
	}

	validatedKinds := []models.Kind{}
	for _, kind := range kinds {
		if kind.Critical && !hasCriticalWrite {
			return []models.Kind{}, postal.UAAScopesError("UAA Scopes Error: Client does not have authority to register critical notifications.")
		}
		kind.ClientID = client.ID
		validatedKinds = append(validatedKinds, kind)
	}

	return validatedKinds, nil
}
