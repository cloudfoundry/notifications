package campaigns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type collectionCreator interface {
	Create(conn collections.ConnectionInterface, campaign collections.Campaign, hasCriticalScope bool) (collections.Campaign, error)
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

	hasCriticalScope := false
	token := context.Get("token").(*jwt.Token)
	for _, scope := range token.Claims["scope"].([]interface{}) {
		if scope.(string) == "critical_notifications.write" {
			hasCriticalScope = true
		}
	}

	database := context.Get("database").(DatabaseInterface)

	campaign, err := h.collection.Create(database.Connection(), collections.Campaign{
		SendTo:         request.SendTo,
		CampaignTypeID: request.CampaignTypeID,
		Text:           request.Text,
		HTML:           request.HTML,
		Subject:        request.Subject,
		TemplateID:     request.TemplateID,
		ReplyTo:        request.ReplyTo,
		ClientID:       context.Get("client_id").(string),
	}, hasCriticalScope)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		case collections.PermissionsError:
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, `{"errors": [%q]}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{ "campaign_id": %q}`, campaign.ID)))
}

func isValid(request createRequest, w http.ResponseWriter, req *http.Request) bool {
	for audienceKey, _ := range request.SendTo {
		if !contains([]string{"user", "space", "org", "email"}, audienceKey) {
			return invalidResponse(w, fmt.Sprintf(`%q is not a valid audience`, audienceKey))
		}

		if audienceKey == "email" {
			if matches := regexp.MustCompile(`[^@]*@{1}[^@]*`).MatchString(request.SendTo[audienceKey]); !matches {
				return invalidResponse(w, fmt.Sprintf(`%q is not a valid email address`, request.SendTo[audienceKey]))
			}
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

func contains(elements []string, element string) bool {
	for _, elem := range elements {
		if element == elem {
			return true
		}
	}

	return false
}

func invalidResponse(w http.ResponseWriter, message string) bool {
	w.WriteHeader(422)
	fmt.Fprintf(w, `{"errors": [%q]}`, message)
	return false
}
