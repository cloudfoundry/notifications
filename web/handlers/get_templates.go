package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type GetTemplates struct {
	Finder      services.TemplateFinderInterface
	ErrorWriter ErrorWriterInterface
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
		Finder:      templateFinder,
		ErrorWriter: errorWriter,
	}
}

func (handler GetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	template, err := handler.Finder.FindByID(context.Get("database").(models.DatabaseInterface), templateID)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(template.Metadata), &metadata)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
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
