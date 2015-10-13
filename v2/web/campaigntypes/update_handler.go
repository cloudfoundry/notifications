package campaigntypes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type collectionUpdater interface {
	Set(conn collections.ConnectionInterface, campaignType collections.CampaignType, clientID string) (collections.CampaignType, error)
	Get(conn collections.ConnectionInterface, campaignTypeID, clientID string) (collections.CampaignType, error)
}

type UpdateHandler struct {
	collection collectionUpdater
}

func NewUpdateHandler(collection collectionUpdater) UpdateHandler {
	return UpdateHandler{
		collection: collection,
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
		validationErrors = append(validationErrors, "name cannot be blank")
	}

	if u.includesDescription() && *u.Description == "" {
		validFlag = false
		validationErrors = append(validationErrors, "description cannot be blank")
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

	updateRequest := UpdateRequest{}

	err := json.NewDecoder(req.Body).Decode(&updateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors": ["invalid json body"]}`))
		return
	}

	validFlag, validationError := updateRequest.isValid()
	if validFlag == false {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"errors": [%q]}`, validationError)
		return
	}

	database := context.Get("database").(DatabaseInterface)
	campaignType, err := h.collection.Get(database.Connection(), campaignTypeID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, `{"errors": [%q]}`, err)
		return
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
			fmt.Fprintf(w, `{ "errors": [%q] }`, "Forbidden: cannot update campaign type with critical flag set to true")
			return
		}
	}

	returnCampaignType, err := h.collection.Set(database.Connection(), campaignType, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, `{"errors": [%q]}`, err)
		return
	}

	json.NewEncoder(w).Encode(NewCampaignTypeResponse(returnCampaignType))
}
