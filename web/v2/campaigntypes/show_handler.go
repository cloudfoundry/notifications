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

	if senderID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, `{"error": %q}`, "missing sender id")
		return
	}

	if campaignTypeID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, `{"error": %q}`, "missing campaign type id")
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(writer, `{"error": %q}`, "missing client id")
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	campaignType, err := h.campaignTypes.Get(database.Connection(), campaignTypeID, senderID, context.Get("client_id").(string))
	if err != nil {
		var errorMessage string
		switch e := err.(type) {
		case collections.NotFoundError:
			errorMessage = e.Message
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
