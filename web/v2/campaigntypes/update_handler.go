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

type UpdateHandler struct {
	campaignTypes collection
}

func NewUpdateHandler(campaignTypes collection) UpdateHandler {
	return UpdateHandler{
		campaignTypes: campaignTypes,
	}
}

type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Critical    *bool   `json:"critical"`
	TemplateID  *string `json:"template_id"`
}

func (u UpdateRequest) isValid() (bool, string) {
	validFlag := true
	validationErrors := make([]string, 0)

	if u.includesName() && *u.Name == "" {
		validFlag = false
		validationErrors = append(validationErrors, "name field cannot be blank")
	}

	if u.includesDescription() && *u.Description == "" {
		validFlag = false
		validationErrors = append(validationErrors, "description field cannot be blank")
	}

	return validFlag, strings.Join(validationErrors, ", ")
}

func (u UpdateRequest) includesName() bool {
	return u.Name != nil
}

func (u UpdateRequest) includesDescription() bool {
	return u.Description != nil
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	campaignTypeID := splitURL[len(splitURL)-1]
	senderID := splitURL[len(splitURL)-3]

	updateRequest := UpdateRequest{}

	err := json.NewDecoder(req.Body).Decode(&updateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid json body"}`))
		return
	}

	validFlag, validationError := updateRequest.isValid()
	if validFlag == false {
		w.WriteHeader(422)
		w.Write([]byte(`{"error": "` + validationError + `"}`))
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	campaignType, err := h.campaignTypes.Get(database.Connection(), senderID, campaignTypeID, context.Get("client_id").(string))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "%s"}`, err.(collections.NotFoundError).Message)
		return
	}

	if updateRequest.includesName() {
		campaignType.Name = *updateRequest.Name
	}

	if updateRequest.includesDescription() {
		campaignType.Description = *updateRequest.Description
	}

	if updateRequest.Critical != nil {
		campaignType.Critical = *updateRequest.Critical
	}

	if updateRequest.TemplateID != nil {
		campaignType.TemplateID = *updateRequest.TemplateID
	}

	returnCampaignType, err := h.campaignTypes.Set(database.Connection(), campaignType, context.Get("client_id").(string))
	if err != nil {
		panic(err)
	}

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
