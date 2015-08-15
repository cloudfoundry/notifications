package templates

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/ryanmoran/stack"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type CreateHandler struct {
	creator     services.TemplateCreatorInterface
	errorWriter errorWriter
}

func NewCreateHandler(creator services.TemplateCreatorInterface, errWriter errorWriter) CreateHandler {
	return CreateHandler{
		creator:     creator,
		errorWriter: errWriter,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateParams, err := NewTemplateParams(req.Body)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	template := templateParams.ToModel()

	templateID, err := h.creator.Create(context.Get("database").(DatabaseInterface), template)
	if err != nil {
		h.errorWriter.Write(w, webutil.TemplateCreateError{})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"template_id":"` + templateID + `"}`))
}

func writeJSON(w http.ResponseWriter, status int, object interface{}) {
	output, err := json.Marshal(object)
	if err != nil {
		panic(err) // No JSON we write into a response should ever panic
	}

	w.WriteHeader(status)
	w.Write(output)
}
