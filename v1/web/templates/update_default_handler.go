package templates

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type UpdateDefaultHandler struct {
	updater     services.TemplateUpdaterInterface
	errorWriter errorWriter
}

func NewUpdateDefaultHandler(updater services.TemplateUpdaterInterface, errWriter errorWriter) UpdateDefaultHandler {
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

	err = h.updater.Update(context.Get("database").(models.DatabaseInterface), models.DefaultTemplateID, template.ToModel())
	if err != nil {
		h.errorWriter.Write(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
