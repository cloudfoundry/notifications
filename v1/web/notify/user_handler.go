package notify

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type UserHandler struct {
	errorWriter errorWriter
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewUserHandler(notify NotifyInterface, errWriter errorWriter, strategy services.StrategyInterface) UserHandler {
	return UserHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h UserHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	conn := context.Get("database").(models.DatabaseInterface).Connection()
	userGUID := strings.TrimPrefix(req.URL.Path, "/users/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(conn, req, context, userGUID, h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
