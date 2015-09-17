package senders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionSetGetter interface {
	collectionSetter
	Get(conn collections.ConnectionInterface, senderID, clientID string) (sender collections.Sender, err error)
}

type UpdateHandler struct {
	senders collectionSetGetter
}

func NewUpdateHandler(senders collectionSetGetter) UpdateHandler {
	return UpdateHandler{
		senders: senders,
	}
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-1]

	var updateRequest struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&updateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors": ["invalid json body"]}`))
		return
	}

	database := context.Get("database").(DatabaseInterface)
	clientID := context.Get("client_id").(string)

	_, err = h.senders.Get(database.Connection(), senderID, clientID)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(fmt.Sprintf(`{ "errors": [ %q ]}`, err)))
		return
	}

	sender, err := h.senders.Set(database.Connection(), collections.Sender{
		ID:       senderID,
		Name:     updateRequest.Name,
		ClientID: clientID,
	})
	if err != nil {
		switch err.(type) {
		case collections.DuplicateRecordError:
			w.WriteHeader(422)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(fmt.Sprintf(`{ "errors": [ %q ]}`, err)))
		return
	}

	json.NewEncoder(w).Encode(NewSenderResponse(sender))
}
