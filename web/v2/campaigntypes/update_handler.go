package campaigntypes

import (
	"encoding/json"
	"net/http"
	"strings"

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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid json body"}`))
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	campaignType, err := h.campaignTypes.Get(database.Connection(), senderID, campaignTypeID, context.Get("client_id").(string))
	// campaignType := collections.CampaignType{
	// 	ID:       campaignTypeID,
	// 	Name:     *updateRequest.Name,
	// 	SenderID: senderID,
	// }

	if updateRequest.Description != nil {
		campaignType.Description = *updateRequest.Description
	}

	if updateRequest.TemplateID != nil {
		campaignType.TemplateID = *updateRequest.TemplateID
	}

	if updateRequest.Critical != nil {
		campaignType.Critical = *updateRequest.Critical
	}

	if updateRequest.Name != nil {
		campaignType.Name = *updateRequest.Name
	}

	returnCampaignType, err := h.campaignTypes.Set(database.Connection(), campaignType, context.Get("client_id").(string))

	jsonMap := map[string]interface{}{
		"id":          returnCampaignType.ID,
		"name":        returnCampaignType.Name,
		"description": returnCampaignType.Description,
		"critical":    returnCampaignType.Critical,
		"template_id": returnCampaignType.TemplateID,
	}

	jsonBody, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(jsonBody))
}
