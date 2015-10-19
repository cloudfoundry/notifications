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

	json.NewEncoder(w).Encode(NewCampaignResponse(campaign))
}
