package senders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionGetter interface {
	Get(conn models.ConnectionInterface, senderID, clientID string) (retrievedSender collections.Sender, err error)
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
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": [ "missing sender id" ] }`))
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "errors": [ "missing client id" ] }`))
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	sender, err := h.senders.Get(database.Connection(), senderID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{ "errors": [ "sender not found" ] }`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "errors": [ "%s" ] }`, err)
		}
		return
	}

	getResponse, _ := json.Marshal(map[string]string{
		"id":   sender.ID,
		"name": sender.Name,
	})

	w.Write(getResponse)
}
