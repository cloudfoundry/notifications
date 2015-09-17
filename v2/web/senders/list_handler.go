package senders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionLister interface {
	List(conn collections.ConnectionInterface, clientID string) ([]collections.Sender, error)
}

type ListHandler struct {
	collection collectionLister
}

func NewListHandler(collection collectionLister) ListHandler {
	return ListHandler{collection: collection}
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	clientID := context.Get("client_id").(string)
	database := context.Get("database").(DatabaseInterface)

	senderList, err := h.collection.List(database.Connection(), clientID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "errors": [ %q ] }`, err)
		return
	}

	json.NewEncoder(w).Encode(NewSendersListResponse(senderList))
}
