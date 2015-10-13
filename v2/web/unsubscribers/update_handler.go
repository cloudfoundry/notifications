package unsubscribers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type UpdateHandler struct {
	collection unsubscribersGetSetter
}

type unsubscribersGetSetter interface {
	Set(connection collections.ConnectionInterface, unsubscriber collections.Unsubscriber) (collections.Unsubscriber, error)
}

func NewUpdateHandler(collection unsubscribersGetSetter) UpdateHandler {
	return UpdateHandler{
		collection: collection,
	}
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	campaignTypeID := splitURL[2]
	userGUID := splitURL[4]

	database := context.Get("database").(DatabaseInterface)
	_, err := h.collection.Set(database.Connection(), collections.Unsubscriber{
		CampaignTypeID: campaignTypeID,
		UserGUID:       userGUID,
	})
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		case collections.PermissionsError:
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(fmt.Sprintf(`{"errors": [%q]}`, err)))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
