package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type ListTemplates struct {
	lister      services.TemplateListerInterface
	errorWriter ErrorWriterInterface
}

func NewListTemplates(templateLister services.TemplateListerInterface, errorWriter ErrorWriterInterface) ListTemplates {
	return ListTemplates{
		lister:      templateLister,
		errorWriter: errorWriter,
	}
}

func (handler ListTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templates, err := handler.lister.List(context.Get("database").(models.DatabaseInterface))
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, templates)
}
