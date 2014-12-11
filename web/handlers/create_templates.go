package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type CreateTemplate struct {
	Creator     services.TemplateCreatorInterface
	ErrorWriter ErrorWriterInterface
}

func NewCreateTemplate(creator services.TemplateCreatorInterface, errorWriter ErrorWriterInterface) CreateTemplate {
	return CreateTemplate{
		Creator:     creator,
		ErrorWriter: errorWriter,
	}
}

func (handler CreateTemplate) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.templates.create",
	}).Log()

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
