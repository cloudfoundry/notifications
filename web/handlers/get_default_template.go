package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type GetDefaultTemplate struct {
	finder      services.TemplateFinderInterface
	errorWriter ErrorWriterInterface
}

func NewGetDefaultTemplate(finder services.TemplateFinderInterface, errorWriter ErrorWriterInterface) GetDefaultTemplate {
	return GetDefaultTemplate{
		finder:      finder,
		errorWriter: errorWriter,
	}
}

func (handler GetDefaultTemplate) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.default_template.get",
	}).Log()

	template, err := handler.finder.FindByID(models.DefaultTemplateID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(template.Metadata), &metadata)
	if err != nil {
		panic(err)
	}

	templateOutput := TemplateOutput{
		Name:     template.Name,
		Subject:  template.Subject,
		HTML:     template.HTML,
		Text:     template.Text,
		Metadata: metadata,
	}

	writeJSON(w, http.StatusOK, templateOutput)
}
