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

type collection interface {
	Add(conn models.ConnectionInterface, campaignType collections.CampaignType, clientID string) (collections.CampaignType, error)
	List(conn models.ConnectionInterface, senderID, clientID string) ([]collections.CampaignType, error)
	Get(conn models.ConnectionInterface, senderID, campaignTypeID, clientID string) (collections.CampaignType, error)
}

type CreateHandler struct {
	campaignTypes collection
}

func NewCreateHandler(campaignTypes collection) CreateHandler {
	return CreateHandler{
		campaignTypes: campaignTypes,
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
		w.Write([]byte(`{"error": "invalid json body"}`))
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
			fmt.Fprintf(w, `{ "error": "%s" }`, http.StatusText(http.StatusForbidden))
			return
		}
	}

	database := context.Get("database").(models.DatabaseInterface)

	campaignType, err := h.campaignTypes.Add(database.Connection(), collections.CampaignType{
		Name:        createRequest.Name,
		Description: createRequest.Description,
		Critical:    createRequest.Critical,
		TemplateID:  createRequest.TemplateID,
		SenderID:    senderID,
	}, context.Get("client_id").(string))
	if err != nil {
		var errorMessage string
		switch e := err.(type) {
		case collections.ValidationError:
			w.WriteHeader(422)
			errorMessage = e.Message
		case collections.NotFoundError:
			w.WriteHeader(404)
			errorMessage = e.Message
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorMessage = err.Error()
		}
		fmt.Fprintf(w, `{ "error": "%s" }`, errorMessage)
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
