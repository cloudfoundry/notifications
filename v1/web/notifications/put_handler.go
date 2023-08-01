package notifications

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanmoran/stack"
)

type registrar interface {
	Register(services.ConnectionInterface, models.Client, []models.Kind) error
	Prune(services.ConnectionInterface, models.Client, []models.Kind) error
}

type PutHandler struct {
	registrar   registrar
	errorWriter errorWriter
}

func NewPutHandler(registrar registrar, errWriter errorWriter) PutHandler {
	return PutHandler{
		registrar:   registrar,
		errorWriter: errWriter,
	}
}

func (h PutHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	database := context.Get("database").(DatabaseInterface)
	connection := database.Connection()

	parameters, err := NewClientRegistrationParams(req.Body)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	err = parameters.Validate()
	if err != nil {
		h.errorWriter.Write(w, err)
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
	claims := token.Claims.(jwt.MapClaims)
	clientID := claims["client_id"].(string)

	client := models.Client{
		ID:          clientID,
		Description: parameters.SourceName,
		TemplateID:  models.DoNotSetTemplateID,
	}

	kinds, err := h.ValidateCriticalScopes(claims["scope"], generatedKinds, client)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	transaction := connection.Transaction()
	transaction.Begin()

	err = h.registrar.Register(transaction, client, kinds)
	if err != nil {
		transaction.Rollback()
		h.errorWriter.Write(w, err)
		return
	}

	if len(parameters.Notifications) > 0 {
		err = h.registrar.Prune(transaction, client, kinds)
		if err != nil {
			transaction.Rollback()
			h.errorWriter.Write(w, err)
			return
		}
	}

	err = transaction.Commit()
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PutHandler) ValidateCriticalScopes(scopes interface{}, kinds []models.Kind, client models.Client) ([]models.Kind, error) {
	hasCriticalWrite := false
	for _, scope := range scopes.([]interface{}) {
		if scope.(string) == "critical_notifications.write" {
			hasCriticalWrite = true
		}
	}

	validatedKinds := []models.Kind{}
	for _, kind := range kinds {
		if kind.Critical && !hasCriticalWrite {
			return []models.Kind{}, webutil.UAAScopesError{Err: errors.New("UAA Scopes Error: Client does not have authority to register critical notifications.")}
		}
		kind.ClientID = client.ID
		validatedKinds = append(validatedKinds, kind)
	}

	return validatedKinds, nil
}
