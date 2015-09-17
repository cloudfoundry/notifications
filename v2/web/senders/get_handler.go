package senders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionGetter interface {
	Get(conn collections.ConnectionInterface, senderID, clientID string) (retrievedSender collections.Sender, err error)
}

type GetHandler struct {
	senders collectionGetter
}

func NewGetHandler(senders collectionGetter) GetHandler {
	return GetHandler{
		senders: senders,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-1]

	if senderID == "" {
		headers := w.Header()
		headers.Set("Location", "/senders")
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "errors": [ "missing client id" ] }`))
		return
	}

	database := context.Get("database").(DatabaseInterface)
	sender, err := h.senders.Get(database.Connection(), senderID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, `{ "errors": [ %q ] }`, err)
		return
	}

	json.NewEncoder(w).Encode(NewSenderResponse(sender))
}
