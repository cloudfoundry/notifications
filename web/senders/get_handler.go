package senders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type GetHandler struct {
	senders collection
}

func NewGetHandler(senders collection) GetHandler {
	return GetHandler{
		senders: senders,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-1]

	database := context.Get("database").(models.DatabaseInterface)
	sender, err := h.senders.Get(database.Connection(), senderID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.ValidationError:
			w.WriteHeader(422)
			w.Write([]byte(`{ "error": "missing sender id" }`))
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{ "error": "sender not found" }`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "error": "%s" }`, err)
		}
		return
	}

	getResponse, _ := json.Marshal(map[string]string{
		"id":   sender.ID,
		"name": sender.Name,
	})

	w.Write(getResponse)
}
