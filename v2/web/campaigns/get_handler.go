package campaigns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type campaignGetter interface {
	Get(connection collections.ConnectionInterface, campaignID, clientID string) (collections.Campaign, error)
}

type GetHandler struct {
	campaigns campaignGetter
}

func NewGetHandler(campaigns campaignGetter) GetHandler {
	return GetHandler{
		campaigns: campaigns,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	campaignID := splitURL[len(splitURL)-1]

	clientID := context.Get("client_id").(string)
	database := context.Get("database").(collections.DatabaseInterface)

	campaign, err := h.campaigns.Get(database.Connection(), campaignID, clientID)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, `{ "errors": [%q] }`, err)
		return
	}

	returnSendTo := map[string][]string{}
	for audienceKey, audience := range campaign.SendTo {
		switch audience.(type) {
		case []interface{}:
			for _, audienceMember := range audience.([]interface{}) {
				returnSendTo[audienceKey] = append(returnSendTo[audienceKey], audienceMember.(string))
			}
		default:
			returnSendTo[audienceKey] = []string{audience.(string)}
		}
	}

	getResponse, _ := json.Marshal(map[string]interface{}{
		"id":               campaign.ID,
		"send_to":          returnSendTo,
		"campaign_type_id": campaign.CampaignTypeID,
		"text":             campaign.Text,
		"html":             campaign.HTML,
		"subject":          campaign.Subject,
		"template_id":      campaign.TemplateID,
		"reply_to":         campaign.ReplyTo,
		"_links": map[string]interface{}{
			"self": map[string]string{
				"href": fmt.Sprintf("/campaigns/%s", campaign.ID),
			},
			"template": map[string]string{
				"href": fmt.Sprintf("/templates/%s", campaign.TemplateID),
			},
			"campaign_type": map[string]string{
				"href": fmt.Sprintf("/campaign_types/%s", campaign.CampaignTypeID),
			},
			"status": map[string]string{
				"href": fmt.Sprintf("/campaigns/%s/status", campaign.ID),
			},
		},
	})
	w.Write(getResponse)
}
