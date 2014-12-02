package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type GetAllNotifications struct {
	finder      services.NotificationsFinderInterface
	errorWriter ErrorWriterInterface
}

func NewGetAllNotifications(notificationsFinder services.NotificationsFinderInterface, errorWriter ErrorWriterInterface) GetAllNotifications {
	return GetAllNotifications{
		finder:      notificationsFinder,
		errorWriter: errorWriter,
	}
}

func (handler GetAllNotifications) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	notifications, err := handler.finder.AllClientNotifications()
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	response, err := json.Marshal(notifications)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}
