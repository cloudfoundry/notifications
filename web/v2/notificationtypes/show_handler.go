package notificationtypes

import (
	"encoding/json"
	"net/http"
	"strings"

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
		panic(err)
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
