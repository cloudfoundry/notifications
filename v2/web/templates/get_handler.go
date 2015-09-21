package templates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionGetter interface {
	Get(conn collections.ConnectionInterface, templateID, clientID string) (collections.Template, error)
}

type GetHandler struct {
	collection collectionGetter
}

func NewGetHandler(collection collectionGetter) GetHandler {
	return GetHandler{
		collection: collection,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	templateID := splitURL[len(splitURL)-1]

	database := context.Get("database").(DatabaseInterface)
	clientID := context.Get("client_id").(string)

	template, err := h.collection.Get(database.Connection(), templateID, clientID)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(fmt.Sprintf(`{"errors": [%q]}`, err)))
		return
	}

	json.NewEncoder(w).Encode(NewTemplateResponse(template))
}
