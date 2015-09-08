package templates

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type templateLister interface {
	List(database services.DatabaseInterface) (templateSummaries map[string]services.TemplateSummary, err error)
}

type ListHandler struct {
	lister      templateLister
	errorWriter errorWriter
}

func NewListHandler(lister templateLister, errWriter errorWriter) ListHandler {
	return ListHandler{
		lister:      lister,
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
