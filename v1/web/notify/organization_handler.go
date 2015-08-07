package notify

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type OrganizationHandler struct {
	errorWriter errorWriter
	notify      NotifyInterface
	strategy    services.StrategyInterface
}

func NewOrganizationHandler(notify NotifyInterface, errWriter errorWriter, strategy services.StrategyInterface) OrganizationHandler {
	return OrganizationHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h OrganizationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	conn := context.Get("database").(models.DatabaseInterface).Connection()
	orgGUID := strings.TrimPrefix(req.URL.Path, "/organizations/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(conn, req, context, orgGUID, h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
