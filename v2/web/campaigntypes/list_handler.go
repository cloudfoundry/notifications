package campaigntypes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionLister interface {
	List(conn collections.ConnectionInterface, senderID, clientID string) ([]collections.CampaignType, error)
}

type ListHandler struct {
	collection collectionLister
}

func NewListHandler(collection collectionLister) ListHandler {
	return ListHandler{
		collection: collection,
	}
}

func (h ListHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) {
	splitURL := strings.Split(request.URL.Path, "/")
	senderID := splitURL[len(splitURL)-2]

	if senderID == "" {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"errors": ["missing sender id"]}`))
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"errors": ["missing client id"]}`))
		return
	}

	database := context.Get("database").(DatabaseInterface)

	campaignTypes, err := h.collection.List(database.Connection(), senderID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(`{"errors": ["sender not found"]}`))
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, `{"errors": [%q]}`, err)
		}
		return
	}

	json.NewEncoder(writer).Encode(NewCampaignTypesListResponse(senderID, campaignTypes))
}
