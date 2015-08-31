package senders

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type sendersDeleter interface {
	Delete(connection collections.ConnectionInterface, senderID, clientID string) error
}

type DeleteHandler struct {
	senders sendersDeleter
}

func NewDeleteHandler(senders sendersDeleter) DeleteHandler {
	return DeleteHandler{
		senders: senders,
	}
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-1]

	conn := context.Get("database").(collections.DatabaseInterface).Connection()
	clientID := context.Get("client_id").(string)

	err := h.senders.Delete(conn, senderID, clientID)
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

	w.WriteHeader(http.StatusNoContent)
}
