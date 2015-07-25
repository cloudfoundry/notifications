package notificationtypes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type ShowHandler struct {
	notificationTypes collection
}

func NewShowHandler(notificationTypes collection) ShowHandler {
	return ShowHandler{
		notificationTypes: notificationTypes,
	}
}

func (h ShowHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) {
	splitURL := strings.Split(request.URL.Path, "/")
	notificationTypeID := splitURL[len(splitURL)-1]
	senderID := splitURL[len(splitURL)-3]

	database := context.Get("database").(models.DatabaseInterface)
	notificationType, err := h.notificationTypes.Get(database.Connection(), notificationTypeID, senderID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.ValidationError:
			writer.WriteHeader(http.StatusBadRequest)
		case collections.NotFoundError:
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		}

		splitErr := strings.Split(err.Error(), ": ")
		fmt.Fprintf(writer, `{"error": "%s"}`, splitErr[len(splitErr)-1])
		return
	}

	jsonMap := map[string]interface{}{
		"id":          notificationType.ID,
		"name":        notificationType.Name,
		"description": notificationType.Description,
		"critical":    notificationType.Critical,
		"template_id": "",
	}

	jsonBody, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	writer.Write([]byte(jsonBody))
}
