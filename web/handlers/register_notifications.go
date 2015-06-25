package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type RegisterNotifications struct {
	registrar   services.RegistrarInterface
	errorWriter ErrorWriterInterface
}

func NewRegisterNotifications(registrar services.RegistrarInterface, errorWriter ErrorWriterInterface) RegisterNotifications {
	return RegisterNotifications{
		registrar:   registrar,
		errorWriter: errorWriter,
	}
}

func (handler RegisterNotifications) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	database := context.Get("database").(models.DatabaseInterface)
	connection := database.Connection()

	parameters, err := params.NewRegistration(req.Body)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	err = parameters.Validate()
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	token := context.Get("token").(*jwt.Token)
	clientID := token.Claims["client_id"].(string)

	client := models.Client{
		ID:          clientID,
		Description: parameters.SourceDescription,
	}

	kinds, err := handler.ValidateCriticalScopes(token.Claims["scope"], parameters.Kinds, client)

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

	if parameters.IncludesKinds {
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
}

func (handler RegisterNotifications) ValidateCriticalScopes(scopes interface{}, kinds []models.Kind, client models.Client) ([]models.Kind, error) {
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
