package templates

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/ryanmoran/stack"
)

type DeleteHandler struct {
	deleter     services.TemplateDeleterInterface
	errorWriter errorWriter
}

func NewDeleteHandler(deleter services.TemplateDeleterInterface, errWriter errorWriter) DeleteHandler {
	return DeleteHandler{
		deleter:     deleter,
		errorWriter: errWriter,
	}
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	err := h.deleter.Delete(context.Get("database").(models.DatabaseInterface), templateID)
	if err != nil {
		h.errorWriter.Write(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
