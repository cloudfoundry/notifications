package senders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionSetter interface {
	Set(conn collections.ConnectionInterface, sender collections.Sender) (createdSender collections.Sender, err error)
}

type CreateHandler struct {
	senders collectionSetter
}

func NewCreateHandler(senders collectionSetter) CreateHandler {
	return CreateHandler{
		senders: senders,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	var createRequest struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "errors": [ "invalid json body" ] }`))
		return
	}

	database := context.Get("database").(DatabaseInterface)

	if createRequest.Name == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": [ "missing sender name" ] }`))
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "errors": [ "missing client id" ] }`))
		return
	}

	sender, err := h.senders.Set(database.Connection(), collections.Sender{
		Name:     createRequest.Name,
		ClientID: context.Get("client_id").(string),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "errors": [ "%s" ] }`, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewSenderResponse(sender))
}
