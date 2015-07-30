package senders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type collection interface {
	Add(conn models.ConnectionInterface, sender collections.Sender) (createdSender collections.Sender, err error)
	Get(conn models.ConnectionInterface, senderID, clientID string) (retrievedSender collections.Sender, err error)
}

type CreateHandler struct {
	senders collection
}

func NewCreateHandler(senders collection) CreateHandler {
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
		w.Write([]byte(`{ "error": "invalid json body" }`))
		return
	}

	database := context.Get("database").(models.DatabaseInterface)

	if createRequest.Name == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "error": "missing sender name" }`))
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "error": "missing client id" }`))
		return
	}

	sender, err := h.senders.Add(database.Connection(), collections.Sender{
		Name:     createRequest.Name,
		ClientID: context.Get("client_id").(string),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "error": "%s" }`, err)
		return
	}

	createResponse, _ := json.Marshal(map[string]string{
		"id":   sender.ID,
		"name": sender.Name,
	})

	w.WriteHeader(http.StatusCreated)
	w.Write(createResponse)
}
