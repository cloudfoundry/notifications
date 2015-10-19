package campaigns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type collectionCreator interface {
	Create(conn collections.ConnectionInterface, campaign collections.Campaign, clientID string, hasCriticalScope bool) (collections.Campaign, error)
}

type clock interface {
	Now() time.Time
}

type CreateHandler struct {
	collection collectionCreator
	clock      clock
}

func NewCreateHandler(collection collectionCreator, clock clock) CreateHandler {
	return CreateHandler{
		collection: collection,
		clock:      clock,
	}
}

type createRequest struct {
	SendTo         map[string][]string `json:"send_to"`
	CampaignTypeID string              `json:"campaign_type_id"`
	Text           string              `json:"text"`
	HTML           string              `json:"html"`
	Subject        string              `json:"subject"`
	TemplateID     string              `json:"template_id"`
	ReplyTo        string              `json:"reply_to"`
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-2]

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
		SenderID:       senderID,
		StartTime:      h.clock.Now(),
	}, context.Get("client_id").(string), hasCriticalScope)
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
	json.NewEncoder(w).Encode(NewCampaignResponse(campaign))
}

func isValid(request createRequest, w http.ResponseWriter, req *http.Request) bool {
	for audienceKey, _ := range request.SendTo {
		if !contains([]string{"users", "spaces", "orgs", "emails"}, audienceKey) {
			return invalidResponse(w, fmt.Sprintf(`%q is not a valid audience`, audienceKey))
		}

		if audienceKey == "emails" {
			if matches := regexp.MustCompile(`[^@]*@{1}[^@]*`).MatchString(request.SendTo[audienceKey][0]); !matches { // TODO: loop over audiences
				return invalidResponse(w, fmt.Sprintf(`%q is not a valid email address`, request.SendTo[audienceKey][0]))
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
