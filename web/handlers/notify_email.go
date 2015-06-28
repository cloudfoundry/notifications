package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type NotifyEmail struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewNotifyEmail(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy services.StrategyInterface) NotifyEmail {
	return NotifyEmail{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (handler NotifyEmail) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)
	database := context.Get("database").(models.DatabaseInterface)

	output, err := handler.notify.Execute(database.Connection(), req, context, "", handler.strategy, params.EmailValidator{}, vcapRequestID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
