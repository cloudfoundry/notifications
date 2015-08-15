package notifications

import (
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type notificationsUpdater interface {
	Update(services.DatabaseInterface, models.Kind) error
}

type UpdateHandler struct {
	updater     notificationsUpdater
	errorWriter errorWriter
}

func NewUpdateHandler(updater notificationsUpdater, errWriter errorWriter) UpdateHandler {
	return UpdateHandler{
		updater:     updater,
		errorWriter: errWriter,
	}
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	var updateParams NotificationUpdateParams

	updateParams, err := NewNotificationParams(req.Body)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	regex := regexp.MustCompile("/clients/(.*)/notifications/(.*)")
	matches := regex.FindStringSubmatch(req.URL.Path)
	clientID, notificationID := matches[1], matches[2]

	err = h.updater.Update(context.Get("database").(DatabaseInterface), updateParams.ToModel(clientID, notificationID))
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
