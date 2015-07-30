package campaigntypes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type UpdateHandler struct {
	campaignTypes collection
}

func NewUpdateHandler(campaignTypes collection) UpdateHandler {
	return UpdateHandler{
		campaignTypes: campaignTypes,
	}
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	campaignTypeID := splitURL[len(splitURL)-1]
	senderID := splitURL[len(splitURL)-3]

	var updateRequest struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Critical    *bool   `json:"critical"`
		TemplateID  *string `json:"template_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&updateRequest)
	if err != nil {
		panic(err)
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte(`{"error": "invalid json body"}`))
		//return
	}

	database := context.Get("database").(models.DatabaseInterface)

	campaignType, err := h.campaignTypes.Set(database.Connection(), collections.CampaignType{
		ID:          campaignTypeID,
		Name:        *updateRequest.Name,
		Description: *updateRequest.Description,
		Critical:    *updateRequest.Critical,
		TemplateID:  *updateRequest.TemplateID,
		SenderID:    senderID,
	}, context.Get("client_id").(string))

	jsonMap := map[string]interface{}{
		"id":          campaignType.ID,
		"name":        campaignType.Name,
		"description": campaignType.Description,
		"critical":    campaignType.Critical,
		"template_id": campaignType.TemplateID,
	}

	jsonBody, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(jsonBody))
}
