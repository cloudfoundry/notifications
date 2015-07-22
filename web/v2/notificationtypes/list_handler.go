package notificationtypes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type ListHandler struct {
	notificationTypes collection
}

func NewListHandler(notificationTypes collection) ListHandler {
	return ListHandler{
		notificationTypes: notificationTypes,
	}
}

func (h ListHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) {
	splitURL := strings.Split(request.URL.Path, "/")
	senderID := splitURL[len(splitURL)-2]

	database := context.Get("database").(models.DatabaseInterface)

	notificationTypes, err := h.notificationTypes.List(database.Connection(), senderID, context.Get("client_id").(string))
	if err != nil {
		panic(err)
	}

	responseList := []interface{}{}

	for _, notificationType := range notificationTypes {
		responseList = append(responseList, map[string]interface{}{
			"id":          notificationType.ID,
			"name":        notificationType.Name,
			"description": notificationType.Description,
			"critical":    notificationType.Critical,
			"template_id": notificationType.TemplateID,
		})
	}

	listResponse, _ := json.Marshal(map[string]interface{}{
		"notification_types": responseList,
	})

	writer.Write(listResponse)
}
