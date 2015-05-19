package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type AssignClientTemplate struct {
	templateAssigner services.TemplateAssignerInterface
	errorWriter      ErrorWriterInterface
}

func NewAssignClientTemplate(assigner services.TemplateAssignerInterface, errorWriter ErrorWriterInterface) AssignClientTemplate {
	return AssignClientTemplate{
		templateAssigner: assigner,
		errorWriter:      errorWriter,
	}
}

type TemplateAssignment struct {
	Template string `json:"template"`
}

func (handler AssignClientTemplate) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	routeRegex := regexp.MustCompile("/clients/(.*)/template")
	clientID := routeRegex.FindStringSubmatch(req.URL.Path)[1]

	var templateAssignment TemplateAssignment
	err := json.NewDecoder(req.Body).Decode(&templateAssignment)
	if err != nil {
		handler.errorWriter.Write(w, params.ParseError{})
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	err = handler.templateAssigner.AssignToClient(database, clientID, templateAssignment.Template)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
