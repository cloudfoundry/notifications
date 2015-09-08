package templates

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type templateUpdater interface {
	Update(database services.DatabaseInterface, templateID string, template models.Template) error
}

type UpdateDefaultHandler struct {
	updater     templateUpdater
	errorWriter errorWriter
}

func NewUpdateDefaultHandler(updater templateUpdater, errWriter errorWriter) UpdateDefaultHandler {
	return UpdateDefaultHandler{
		updater:     updater,
		errorWriter: errWriter,
	}
}

func (h UpdateDefaultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	template, err := NewTemplateParams(req.Body)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	err = h.updater.Update(context.Get("database").(DatabaseInterface), models.DefaultTemplateID, template.ToModel())
	if err != nil {
		h.errorWriter.Write(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
