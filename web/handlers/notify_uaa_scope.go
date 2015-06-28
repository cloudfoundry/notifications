package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type NotifyUAAScope struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewNotifyUAAScope(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy services.StrategyInterface) NotifyUAAScope {
	return NotifyUAAScope{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (handler NotifyUAAScope) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := context.Get("database").(models.DatabaseInterface).Connection()
	scope := strings.TrimPrefix(req.URL.Path, "/uaa_scopes/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := handler.notify.Execute(connection, req, context, scope, handler.strategy, params.GUIDValidator{}, vcapRequestID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
