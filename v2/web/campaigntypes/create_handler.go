package campaigntypes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type collectionSetter interface {
	Set(conn models.ConnectionInterface, campaignType collections.CampaignType, clientID string) (collections.CampaignType, error)
}

type CreateHandler struct {
	collection collectionSetter
}

func NewCreateHandler(collection collectionSetter) CreateHandler {
	return CreateHandler{
		collection: collection,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-2]

	var createRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Critical    bool   `json:"critical"`
		TemplateID  string `json:"template_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors": ["invalid json body"]}`))
		return
	}

	if createRequest.Name == "" {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{ "errors": [%q]}`, "missing campaign type name")
		return
	}

	if createRequest.Description == "" {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"errors": [%q]}`, "missing campaign type description")
		return
	}

	if createRequest.Critical == true {
		hasCriticalWrite := false
		token := context.Get("token").(*jwt.Token)
		for _, scope := range token.Claims["scope"].([]interface{}) {
			if scope.(string) == "critical_notifications.write" {
				hasCriticalWrite = true
			}
		}

		if hasCriticalWrite == false {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{ "errors": [%q]}`, http.StatusText(http.StatusForbidden))
			return
		}
	}

	database := context.Get("database").(models.DatabaseInterface)

	campaignType, err := h.collection.Set(database.Connection(), collections.CampaignType{
		Name:        createRequest.Name,
		Description: createRequest.Description,
		Critical:    createRequest.Critical,
		TemplateID:  createRequest.TemplateID,
		SenderID:    senderID,
	}, context.Get("client_id").(string))
	if err != nil {
		var errorMessage string
		switch e := err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
			errorMessage = e.Message
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorMessage = err.Error()
		}
		fmt.Fprintf(w, `{ "errors": [%q]}`, errorMessage)
		return
	}

	createResponse, _ := json.Marshal(map[string]interface{}{
		"id":          campaignType.ID,
		"name":        campaignType.Name,
		"description": campaignType.Description,
		"critical":    campaignType.Critical,
		"template_id": campaignType.TemplateID,
	})

	w.WriteHeader(http.StatusCreated)
	w.Write(createResponse)
}
