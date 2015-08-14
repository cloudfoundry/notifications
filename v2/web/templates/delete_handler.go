package templates

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionsDeleter interface {
	Get(conn collections.ConnectionInterface, templateID, clientID string) (collections.Template, error)
	Delete(conn collections.ConnectionInterface, templateID string) error
}

type DeleteHandler struct {
	deleter collectionsDeleter
}

func NewDeleteHandler(deleter collectionsDeleter) DeleteHandler {
	return DeleteHandler{
		deleter: deleter,
	}
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, request *http.Request, context stack.Context) {
	splitURL := strings.Split(request.URL.Path, "/")
	templateID := splitURL[len(splitURL)-1]

	database := context.Get("database").(collections.DatabaseInterface)
	clientID := context.Get("client_id").(string)

	_, err := h.deleter.Get(database.Connection(), templateID, clientID)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, `{"errors": [%q]}`, err)
		return
	}

	err = h.deleter.Delete(database.Connection(), templateID)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, `{"errors": [%q]}`, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
