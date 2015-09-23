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

type collectionSetter interface {
	Set(conn collections.ConnectionInterface, campaignType collections.CampaignType, clientID string) (collections.CampaignType, error)
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
			w.Write([]byte(`{ "errors": [ "You do not have permission to create critical campaign types" ]}`))
			return
		}
	}

	database := context.Get("database").(DatabaseInterface)

	campaignType, err := h.collection.Set(database.Connection(), collections.CampaignType{
		Name:        createRequest.Name,
		Description: createRequest.Description,
		Critical:    createRequest.Critical,
		TemplateID:  createRequest.TemplateID,
		SenderID:    senderID,
	}, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, `{ "errors": [%q]}`, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewCampaignTypeResponse(campaignType))
}
