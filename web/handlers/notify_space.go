package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/ryanmoran/stack"
)

type NotifySpace struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewNotifySpace(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy services.StrategyInterface) NotifySpace {
	return NotifySpace{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := context.Get("database").(models.DatabaseInterface).Connection()
	spaceGUID := strings.TrimPrefix(req.URL.Path, "/spaces/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := handler.notify.Execute(connection, req, context, spaceGUID, handler.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
