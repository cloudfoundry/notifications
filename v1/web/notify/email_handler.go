package notify

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

const (
	VCAPRequestIDKey    = "vcap_request_id"
	RequestReceivedTime = "request_received_time"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type EmailHandler struct {
	errorWriter errorWriter
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewEmailHandler(notify NotifyInterface, errWriter errorWriter, strategy services.StrategyInterface) EmailHandler {
	return EmailHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h EmailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)
	database := context.Get("database").(db.DatabaseInterface)
	conn := database.Connection()

	output, err := h.notify.Execute(conn, req, context, "", h.strategy, EmailValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
