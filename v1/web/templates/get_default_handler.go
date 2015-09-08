package templates

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type templateFinder interface {
	FindByID(database services.DatabaseInterface, templateID string) (models.Template, error)
}

type GetDefaultHandler struct {
	finder      templateFinder
	errorWriter errorWriter
}

func NewGetDefaultHandler(finder templateFinder, errWriter errorWriter) GetDefaultHandler {
	return GetDefaultHandler{
		finder:      finder,
		errorWriter: errWriter,
	}
}

func (h GetDefaultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	template, err := h.finder.FindByID(context.Get("database").(DatabaseInterface), models.DefaultTemplateID)
	if err != nil {
		h.errorWriter.Write(w, err)
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
