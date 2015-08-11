package notify

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type SpaceHandler struct {
	errorWriter errorWriter
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewSpaceHandler(notify NotifyInterface, errWriter errorWriter, strategy services.StrategyInterface) SpaceHandler {
	return SpaceHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h SpaceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	conn := context.Get("database").(db.DatabaseInterface).Connection()
	spaceGUID := strings.TrimPrefix(req.URL.Path, "/spaces/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(conn, req, context, spaceGUID, h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
