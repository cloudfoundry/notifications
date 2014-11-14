package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type UnsetTemplates struct {
	deleter     services.TemplateDeleterInterface
	errorWriter ErrorWriterInterface
}

func NewUnsetTemplates(deleter services.TemplateDeleterInterface, errorWriter ErrorWriterInterface) UnsetTemplates {
	return UnsetTemplates{
		deleter:     deleter,
		errorWriter: errorWriter,
	}
}

func (handler UnsetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, stack stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.templates.delete",
	}).Log()

	templateName := strings.Split(req.URL.Path, "/templates/")[1]

	err := handler.deleter.Delete(templateName)
	if err != nil {
		handler.errorWriter.Write(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
