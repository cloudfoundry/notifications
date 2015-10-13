package campaigntypes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionGetter interface {
	Get(conn collections.ConnectionInterface, campaignTypeID, clientID string) (collections.CampaignType, error)
}

type ShowHandler struct {
	collection collectionGetter
}

func NewShowHandler(collection collectionGetter) ShowHandler {
	return ShowHandler{
		collection: collection,
	}
}

func (h ShowHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) {
	splitURL := strings.Split(request.URL.Path, "/")
	campaignTypeID := splitURL[len(splitURL)-1]

	clientID := context.Get("client_id")
	if clientID == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(writer, `{"errors": [%q]}`, "missing client id")
		return
	}

	database := context.Get("database").(DatabaseInterface)
	campaignType, err := h.collection.Get(database.Connection(), campaignTypeID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(writer, `{"errors": [%q]}`, err)
		return
	}

	json.NewEncoder(writer).Encode(NewCampaignTypeResponse(campaignType))
}
