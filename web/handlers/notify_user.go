package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/ryanmoran/stack"
)

type NotifyUser struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewNotifyUser(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy services.StrategyInterface) NotifyUser {
	return NotifyUser{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := context.Get("database").(models.DatabaseInterface).Connection()

	userGUID := strings.TrimPrefix(req.URL.Path, "/users/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := handler.notify.Execute(connection, req, context, userGUID, handler.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
