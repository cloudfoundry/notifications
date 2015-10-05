package clients

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/v1/collections"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/ryanmoran/stack"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type assignsTemplates interface {
	AssignToClient(connection collections.ConnectionInterface, clientID, templateID string) error
}

type AssignTemplateHandler struct {
	templateAssigner assignsTemplates
	errorWriter      errorWriter
}

func NewAssignTemplateHandler(assigner assignsTemplates, errWriter errorWriter) AssignTemplateHandler {
	return AssignTemplateHandler{
		templateAssigner: assigner,
		errorWriter:      errWriter,
	}
}

type TemplateAssignment struct {
	Template string `json:"template"`
}

func (h AssignTemplateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	routeRegex := regexp.MustCompile("/clients/(.*)/template")
	clientID := routeRegex.FindStringSubmatch(req.URL.Path)[1]

	var templateAssignment TemplateAssignment
	err := json.NewDecoder(req.Body).Decode(&templateAssignment)
	if err != nil {
		h.errorWriter.Write(w, webutil.ParseError{})
		return
	}

	database := context.Get("database").(DatabaseInterface)
	err = h.templateAssigner.AssignToClient(database.Connection(), clientID, templateAssignment.Template)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
