package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type CreateTemplate struct {
	creator     services.TemplateCreatorInterface
	errorWriter ErrorWriterInterface
}

func NewCreateTemplate(creator services.TemplateCreatorInterface, errorWriter ErrorWriterInterface) CreateTemplate {
	return CreateTemplate{
		creator:     creator,
		errorWriter: errorWriter,
	}
}

func (handler CreateTemplate) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateParams, err := params.NewTemplateParams(req.Body)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	template := templateParams.ToModel()

	templateID, err := handler.creator.Create(context.Get("database").(models.DatabaseInterface), template)
	if err != nil {
		handler.errorWriter.Write(w, params.TemplateCreateError{})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"template_id":"` + templateID + `"}`))
}
