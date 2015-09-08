package templates

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ryanmoran/stack"
)

type TemplateOutput struct {
	Name     string                 `json:"name"`
	Subject  string                 `json:"subject"`
	HTML     string                 `json:"html"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata"`
}

type GetHandler struct {
	finder      templateFinder
	errorWriter errorWriter
}

func NewGetHandler(templateFinder templateFinder, errWriter errorWriter) GetHandler {
	return GetHandler{
		finder:      templateFinder,
		errorWriter: errWriter,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	template, err := h.finder.FindByID(context.Get("database").(DatabaseInterface), templateID)
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
