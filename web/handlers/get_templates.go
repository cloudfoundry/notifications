package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type GetTemplates struct {
	Finder      services.TemplateFinderInterface
	ErrorWriter ErrorWriterInterface
}

type TemplateOutput struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
}

func NewGetTemplates(templateFinder services.TemplateFinderInterface, errorWriter ErrorWriterInterface) GetTemplates {
	return GetTemplates{
		Finder:      templateFinder,
		ErrorWriter: errorWriter,
	}
}

func (handler GetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.templates.get",
	}).Log()

	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	template, err := handler.Finder.Find(templateID)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	templateOutput := TemplateOutput{
		Name:    template.Name,
		Subject: template.Subject,
		HTML:    template.HTML,
		Text:    template.Text,
	}

	response, err := json.Marshal(templateOutput)
	if err != nil {
		panic(err)
	}
	w.Write(response)
}
