package templates

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type GetHandler struct {
	finder      services.TemplateFinderInterface
	errorWriter errorWriter
}

type TemplateOutput struct {
	Name     string                 `json:"name"`
	Subject  string                 `json:"subject"`
	HTML     string                 `json:"html"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata"`
}

func NewGetHandler(templateFinder services.TemplateFinderInterface, errWriter errorWriter) GetHandler {
	return GetHandler{
		finder:      templateFinder,
		errorWriter: errWriter,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	template, err := h.finder.FindByID(context.Get("database").(models.DatabaseInterface), templateID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(template.Metadata), &metadata)
	if err != nil {
		h.errorWriter.Write(w, err)
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
