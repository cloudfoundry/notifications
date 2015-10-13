package campaigntypes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type DeleteHandler struct {
	collection collectionDeleter
}

type collectionDeleter interface {
	Delete(conn collections.ConnectionInterface, campaignTypeID, clientID string) error
}

func NewDeleteHandler(campaignTypesCollection collectionDeleter) DeleteHandler {
	return DeleteHandler{
		collection: campaignTypesCollection,
	}
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	campaignTypeID := splitURL[len(splitURL)-1]

	database := context.Get("database").(DatabaseInterface)
	if err := h.collection.Delete(database.Connection(), campaignTypeID, context.Get("client_id").(string)); err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"errors": [%q]}`, err)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"errors": ["Delete failed with error: %s"]}`, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
