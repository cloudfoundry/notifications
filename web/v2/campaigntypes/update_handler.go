package campaigntypes

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
		validationErrors = append(validationErrors, "name can not be blank")
	}

	if u.includesDescription() && *u.Description == "" {
		validFlag = false
		validationErrors = append(validationErrors, "description can not be blank")
	}

	return validFlag, strings.Join(validationErrors, ", ")
}

func (u UpdateRequest) includesName() bool {
	return u.Name != nil
}

func (u UpdateRequest) includesDescription() bool {
	return u.Description != nil
}

func (u UpdateRequest) includesCritical() bool {
	return u.Critical != nil
}

func (u UpdateRequest) includesTemplateID() bool {
	return u.TemplateID != nil
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
		fmt.Fprintf(w, `{"error": %q}`, validationError)
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	campaignType, err := h.campaignTypes.Get(database.Connection(), campaignTypeID, senderID, context.Get("client_id").(string))
	if err != nil {
		switch err := err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error": %q}`, err.Message)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": %q}`, err.Error())
			return
		}
	}

	if updateRequest.includesName() {
		campaignType.Name = *updateRequest.Name
	}

	if updateRequest.includesDescription() {
		campaignType.Description = *updateRequest.Description
	}

	if updateRequest.includesCritical() {
		campaignType.Critical = *updateRequest.Critical
	}

	if updateRequest.includesTemplateID() {
		campaignType.TemplateID = *updateRequest.TemplateID
	}

	if campaignType.Critical == true {
		hasCriticalWrite := false
		token := context.Get("token").(*jwt.Token)
		for _, scope := range token.Claims["scope"].([]interface{}) {
			if scope.(string) == "critical_notifications.write" {
				hasCriticalWrite = true
			}
		}

		if hasCriticalWrite == false {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{ "error": %q }`, "Forbidden: can not update campaign type with critical flag set to true")
			return
		}
	}

	returnCampaignType, err := h.campaignTypes.Set(database.Connection(), campaignType, context.Get("client_id").(string))
	if err != nil {
		switch err := err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error": %q}`, err.Message)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": %q}`, err.Error())
			return
		}
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
