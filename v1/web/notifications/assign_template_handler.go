package notifications

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/ryanmoran/stack"
)

type TemplateAssignment struct {
	Template string `json:"template"`
}

type AssignTemplateHandler struct {
	templateAssigner services.TemplateAssignerInterface
	errorWriter      errorWriter
}

func NewAssignTemplateHandler(assigner services.TemplateAssignerInterface, errWriter errorWriter) AssignTemplateHandler {
	return AssignTemplateHandler{
		templateAssigner: assigner,
		errorWriter:      errWriter,
	}
}

func (h AssignTemplateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	clientID, notificationID := h.parseURL(req.URL.Path)

	var templateAssignment TemplateAssignment
	err := json.NewDecoder(req.Body).Decode(&templateAssignment)
	if err != nil {
		h.errorWriter.Write(w, webutil.ParseError{})
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	err = h.templateAssigner.AssignToNotification(database, clientID, notificationID, templateAssignment.Template)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h AssignTemplateHandler) parseURL(path string) (string, string) {
	routeMatches := regexp.MustCompile("/clients/(.*)/notifications/(.*)/template").FindStringSubmatch(path)

	return routeMatches[1], routeMatches[2]
}
