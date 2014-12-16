package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type AssignNotificationTemplate struct {
	templateAssigner services.TemplateAssignerInterface
	errorWriter      ErrorWriterInterface
}

func NewAssignNotificationTemplate(assigner services.TemplateAssignerInterface, errorWriter ErrorWriterInterface) AssignNotificationTemplate {
	return AssignNotificationTemplate{
		templateAssigner: assigner,
		errorWriter:      errorWriter,
	}
}

func (handler AssignNotificationTemplate) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	clientID, notificationID := handler.parseURL(req.URL.Path)

	var templateAssignment TemplateAssignment
	err := json.NewDecoder(req.Body).Decode(&templateAssignment)
	if err != nil {
		handler.errorWriter.Write(w, params.ParseError{})
		return
	}

	err = handler.templateAssigner.AssignToNotification(clientID, notificationID, templateAssignment.Template)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler AssignNotificationTemplate) parseURL(path string) (string, string) {
	routeMatches := regexp.MustCompile("/clients/(.*)/notifications/(.*)/template").FindStringSubmatch(path)

	return routeMatches[1], routeMatches[2]
}
