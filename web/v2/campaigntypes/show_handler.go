package campaigntypes

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
	campaignTypes collection
}

func NewShowHandler(campaignTypes collection) ShowHandler {
	return ShowHandler{
		campaignTypes: campaignTypes,
	}
}

func (h ShowHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) {
	splitURL := strings.Split(request.URL.Path, "/")
	campaignTypeID := splitURL[len(splitURL)-1]
	senderID := splitURL[len(splitURL)-3]

	database := context.Get("database").(models.DatabaseInterface)
	campaignType, err := h.campaignTypes.Get(database.Connection(), campaignTypeID, senderID, context.Get("client_id").(string))
	if err != nil {
		var errorMessage string
		switch err.(type) {
		case collections.ValidationError:
			errorMessage = err.(collections.ValidationError).Message
			writer.WriteHeader(http.StatusBadRequest)
		case collections.NotFoundError:
			errorMessage = err.(collections.NotFoundError).Message
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			errorMessage = err.Error()
		}

		fmt.Fprintf(writer, `{"error": "%s"}`, errorMessage)
		return
	}

	jsonMap := map[string]interface{}{
		"id":          campaignType.ID,
		"name":        campaignType.Name,
		"description": campaignType.Description,
		"critical":    campaignType.Critical,
		"template_id": "",
	}

	jsonBody, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	writer.Write([]byte(jsonBody))
}
