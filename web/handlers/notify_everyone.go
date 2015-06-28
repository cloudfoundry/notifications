package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type NotifyEveryone struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewNotifyEveryone(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy services.StrategyInterface) NotifyEveryone {
	return NotifyEveryone{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (handler NotifyEveryone) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := context.Get("database").(models.DatabaseInterface).Connection()
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := handler.notify.Execute(connection, req, context, "", handler.strategy, params.GUIDValidator{}, vcapRequestID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
