package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type UpdateDefaultTemplate struct {
	updater     services.TemplateUpdaterInterface
	errorWriter ErrorWriterInterface
}

func NewUpdateDefaultTemplate(updater services.TemplateUpdaterInterface, errorWriter ErrorWriterInterface) UpdateDefaultTemplate {
	return UpdateDefaultTemplate{
		updater:     updater,
		errorWriter: errorWriter,
	}
}

func (handler UpdateDefaultTemplate) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	template, err := params.NewTemplate(req.Body)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	err = handler.updater.Update(models.DefaultTemplateID, template.ToModel())
	if err != nil {
		handler.errorWriter.Write(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
