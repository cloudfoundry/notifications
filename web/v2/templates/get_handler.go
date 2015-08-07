package templates

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type collectionGetter interface {
	Get(conn models.ConnectionInterface, templateID, clientID string) (collections.Template, error)
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

	database := context.Get("database").(models.DatabaseInterface)
	clientID := context.Get("client_id").(string)

	template, _ := h.collection.Get(database.Connection(), templateID, clientID)

	json, err := json.Marshal(template)
	if err != nil {
		panic(err)
	}
	w.Write(json)
}
