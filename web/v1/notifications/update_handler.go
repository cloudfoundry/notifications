package notifications

import (
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type UpdateHandler struct {
	updater     NotificationsUpdaterInterface
	errorWriter errorWriter
}

func NewUpdateHandler(notificationsUpdater NotificationsUpdaterInterface, errWriter errorWriter) UpdateHandler {
	return UpdateHandler{
		updater:     notificationsUpdater,
		errorWriter: errWriter,
	}
}

type NotificationsUpdaterInterface interface {
	Update(models.DatabaseInterface, models.Kind) error
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

	err = h.updater.Update(context.Get("database").(models.DatabaseInterface), updateParams.ToModel(clientID, notificationID))
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
