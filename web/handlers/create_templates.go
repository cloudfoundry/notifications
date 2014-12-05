package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type CreateTemplates struct {
	Creator     services.TemplateCreatorInterface
	ErrorWriter ErrorWriterInterface
}

func NewCreateTemplates(creator services.TemplateCreatorInterface, errorWriter ErrorWriterInterface) CreateTemplates {
	return CreateTemplates{
		Creator:     creator,
		ErrorWriter: errorWriter,
	}
}

func (handler CreateTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateParams, err := params.NewTemplate(req.Body)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}
	templateID, err := handler.Creator.Create(templateParams.ToModel())
	if err != nil {
		handler.ErrorWriter.Write(w, params.TemplateCreateError{})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"template-id":"` + templateID + `"}`))

}
