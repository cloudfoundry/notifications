package notificationtypes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type collection interface {
	Add(conn models.ConnectionInterface, notificationType collections.NotificationType, clientID string) (createdNotificationType collections.NotificationType, err error)
	List(conn models.ConnectionInterface, senderID, clientID string) (notificationTypes []collections.NotificationType, err error)
	Get(conn models.ConnectionInterface, senderID, notificationTypeID, clientID string) (notificationType collections.NotificationType, err error)
}

type CreateHandler struct {
	notificationTypes collection
}

func NewCreateHandler(notificationTypes collection) CreateHandler {
	return CreateHandler{
		notificationTypes: notificationTypes,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-2]

	var createRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Critical    bool   `json:"critical"`
		TemplateID  string `json:"template_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid json body"}`))
		return
	}

	if createRequest.Critical == true {
		hasCriticalWrite := false
		token := context.Get("token").(*jwt.Token)
		for _, scope := range token.Claims["scope"].([]interface{}) {
			if scope.(string) == "critical_notifications.write" {
				hasCriticalWrite = true
			}
		}

		if hasCriticalWrite == false {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{ "error": "%s" }`, http.StatusText(http.StatusForbidden))
			return
		}
	}

	database := context.Get("database").(models.DatabaseInterface)

	notificationType, err := h.notificationTypes.Add(database.Connection(), collections.NotificationType{
		Name:        createRequest.Name,
		Description: createRequest.Description,
		Critical:    createRequest.Critical,
		TemplateID:  createRequest.TemplateID,
		SenderID:    senderID,
	}, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.ValidationError:
			w.WriteHeader(422)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, `{ "error": "%s" }`, err)
		return
	}

	createResponse, _ := json.Marshal(map[string]interface{}{
		"id":          notificationType.ID,
		"name":        notificationType.Name,
		"description": notificationType.Description,
		"critical":    notificationType.Critical,
		"template_id": notificationType.TemplateID,
	})

	w.WriteHeader(http.StatusCreated)
	w.Write(createResponse)
}
