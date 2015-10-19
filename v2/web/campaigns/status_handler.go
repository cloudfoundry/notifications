package campaigns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type campaignStatusGetter interface {
	Get(connection collections.ConnectionInterface, campaignID, senderID string) (collections.CampaignStatus, error)
}

type StatusHandler struct {
	campaignStatuses campaignStatusGetter
}

func NewStatusHandler(campaignStatuses campaignStatusGetter) StatusHandler {
	return StatusHandler{
		campaignStatuses: campaignStatuses,
	}
}

func (h StatusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	campaignID := splitURL[len(splitURL)-2]
	clientID := context.Get("client_id").(string)
	conn := context.Get("database").(collections.DatabaseInterface).Connection()

	status, err := h.campaignStatuses.Get(conn, campaignID, clientID)
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

	json.NewEncoder(w).Encode(NewCampaignStatusResponse(status))
}
