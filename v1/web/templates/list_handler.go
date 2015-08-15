package templates

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type ListHandler struct {
	lister      services.TemplateListerInterface
	errorWriter errorWriter
}

func NewListHandler(templateLister services.TemplateListerInterface, errWriter errorWriter) ListHandler {
	return ListHandler{
		lister:      templateLister,
		errorWriter: errWriter,
	}
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templates, err := h.lister.List(context.Get("database").(DatabaseInterface))
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, templates)
}
