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

type ListHandler struct {
	campaignTypes collection
}

func NewListHandler(campaignTypes collection) ListHandler {
	return ListHandler{
		campaignTypes: campaignTypes,
	}
}

func (h ListHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) {
	splitURL := strings.Split(request.URL.Path, "/")
	senderID := splitURL[len(splitURL)-2]

	if senderID == "" {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{ "error": "missing sender id" }`))
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{ "error": "missing client id" }`))
		return
	}

	database := context.Get("database").(models.DatabaseInterface)

	campaignTypes, err := h.campaignTypes.List(database.Connection(), senderID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(`{ "error": "sender not found" }`))
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, `{ "error": "%s" }`, err)
		}
		return
	}

	responseList := []interface{}{}

	for _, campaignType := range campaignTypes {
		responseList = append(responseList, map[string]interface{}{
			"id":          campaignType.ID,
			"name":        campaignType.Name,
			"description": campaignType.Description,
			"critical":    campaignType.Critical,
			"template_id": campaignType.TemplateID,
		})
	}

	listResponse, _ := json.Marshal(map[string]interface{}{
		"campaign_types": responseList,
	})

	writer.Write(listResponse)
}
