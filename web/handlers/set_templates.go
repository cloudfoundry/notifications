package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type SetTemplates struct {
	updater     services.TemplateUpdaterInterface
	ErrorWriter ErrorWriterInterface
}

func NewSetTemplates(updater services.TemplateUpdaterInterface, errorWriter ErrorWriterInterface) SetTemplates {
	return SetTemplates{
		updater:     updater,
		ErrorWriter: errorWriter,
	}
}

func (handler SetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.templates.put",
	}).Log()

	templateName := strings.Split(req.URL.String(), "/templates/")[1]

	templateParams, err := params.NewTemplate(templateName, req.Body)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	err = templateParams.Validate()
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	err = handler.updater.Update(templateParams.ToModel())
	if err != nil {
		handler.ErrorWriter.Write(w, params.TemplateUpdateError{})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
