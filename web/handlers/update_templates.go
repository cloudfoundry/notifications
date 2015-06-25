package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type UpdateTemplates struct {
	updater     services.TemplateUpdaterInterface
	ErrorWriter ErrorWriterInterface
}

func NewUpdateTemplates(updater services.TemplateUpdaterInterface, errorWriter ErrorWriterInterface) UpdateTemplates {
	return UpdateTemplates{
		updater:     updater,
		ErrorWriter: errorWriter,
	}
}

func (handler UpdateTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := strings.Split(req.URL.String(), "/templates/")[1]

	templateParams, err := params.NewTemplate(req.Body)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	err = handler.updater.Update(context.Get("database").(models.DatabaseInterface), templateID, templateParams.ToModel())
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
