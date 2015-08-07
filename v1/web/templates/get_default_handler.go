package templates

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/ryanmoran/stack"
)

type GetDefaultHandler struct {
	finder      services.TemplateFinderInterface
	errorWriter errorWriter
}

func NewGetDefaultHandler(finder services.TemplateFinderInterface, errWriter errorWriter) GetDefaultHandler {
	return GetDefaultHandler{
		finder:      finder,
		errorWriter: errWriter,
	}
}

func (h GetDefaultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	template, err := h.finder.FindByID(context.Get("database").(models.DatabaseInterface), models.DefaultTemplateID)
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
