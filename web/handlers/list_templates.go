package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type ListTemplates struct {
	Lister      services.TemplateListerInterface
	ErrorWriter ErrorWriterInterface
}

func NewListTemplates(templateLister services.TemplateListerInterface, errorWriter ErrorWriterInterface) ListTemplates {
	return ListTemplates{
		Lister:      templateLister,
		ErrorWriter: errorWriter,
	}
}

func (handler ListTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.templates.list",
	}).Log()

	templates, err := handler.Lister.List()
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, templates)
}
