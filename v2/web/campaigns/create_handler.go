package campaigns

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionCreator interface {
	Create(conn collections.ConnectionInterface, campaign collections.Campaign) (collections.Campaign, error)
}

type CreateHandler struct {
	collection collectionCreator
}

func NewCreateHandler(collection collectionCreator) CreateHandler {
	return CreateHandler{
		collection: collection,
	}
}

type createRequest struct {
	SendTo         map[string]string `json:"send_to"`
	CampaignTypeID string            `json:"campaign_type_id"`
	Text           string            `json:"text"`
	HTML           string            `json:"html"`
	Subject        string            `json:"subject"`
	TemplateID     string            `json:"template_id"`
	ReplyTo        string            `json:"reply_to"`
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	var request createRequest

	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors": [%q]}`, "invalid json body")
		return
	}

	if !isValid(request, w, req) {
		return
	}

	database := context.Get("database").(DatabaseInterface)

	campaign, err := h.collection.Create(database.Connection(), collections.Campaign{
		SendTo:         request.SendTo,
		CampaignTypeID: request.CampaignTypeID,
		Text:           request.Text,
		HTML:           request.HTML,
		Subject:        request.Subject,
		TemplateID:     request.TemplateID,
		ReplyTo:  request.ReplyTo,
		ClientID: context.Get("client_id").(string),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors": [%q]}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{ "campaign_id": %q}`, campaign.ID)))
}

func isValid(request createRequest, w http.ResponseWriter, req *http.Request) bool {
	for audienceKey, _ := range request.SendTo {
		if audienceKey != "user" {
			return invalidResponse(w, fmt.Sprintf(`%q is not a valid audience`, audienceKey))
		}
	}

	if request.CampaignTypeID == "" {
		return invalidResponse(w, "missing campaign_type_id")
	}

	if request.Text == "" && request.HTML == "" {
		return invalidResponse(w, "missing either campaign text or html")
	}

	if request.Subject == "" {
		return invalidResponse(w, "missing subject")
	}

	return true
}

func invalidResponse(w http.ResponseWriter, message string) bool {
	w.WriteHeader(422)
	fmt.Fprintf(w, `{"errors": [%q]}`, message)
	return false
}
