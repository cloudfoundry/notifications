package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/ryanmoran/stack"
)

type GetTemplates struct {
	finder      services.TemplateFinderInterface
	errorWriter ErrorWriterInterface
}

type TemplateOutput struct {
	Name     string                 `json:"name"`
	Subject  string                 `json:"subject"`
	HTML     string                 `json:"html"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata"`
}

func NewGetTemplates(templateFinder services.TemplateFinderInterface, errorWriter ErrorWriterInterface) GetTemplates {
	return GetTemplates{
		finder:      templateFinder,
		errorWriter: errorWriter,
	}
}

func (handler GetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	template, err := handler.finder.FindByID(context.Get("database").(models.DatabaseInterface), templateID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(template.Metadata), &metadata)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
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
