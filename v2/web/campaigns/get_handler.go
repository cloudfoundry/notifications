package campaigns

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type campaignGetter interface {
	Get(connection collections.ConnectionInterface, campaignID, senderID, clientID string) (collections.Campaign, error)
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
	senderID := splitURL[len(splitURL)-3]
	campaignID := splitURL[len(splitURL)-1]

	clientID := context.Get("client_id").(string)
	database := context.Get("database").(db.DatabaseInterface)

	campaign, err := h.campaigns.Get(database.Connection(), campaignID, senderID, clientID)
	if err != nil {
		panic(err)
	}
	getResponse, _ := json.Marshal(map[string]interface{}{
		"send_to":          campaign.SendTo,
		"campaign_type_id": campaign.CampaignTypeID,
		"text":             campaign.Text,
		"html":             campaign.HTML,
		"subject":          campaign.Subject,
		"template_id":      campaign.TemplateID,
		"reply_to":         campaign.ReplyTo,
	})
	w.Write(getResponse)
}
