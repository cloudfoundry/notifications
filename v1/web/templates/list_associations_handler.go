package templates

import (
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type TemplateAssociation struct {
	Client       string `json:"client"`
	Notification string `json:"notification,omitempty"`
}

type templateAssociationLister interface {
	List(database services.DatabaseInterface, templateID string) ([]services.TemplateAssociation, error)
}

type ListAssociationsHandler struct {
	lister      templateAssociationLister
	errorWriter errorWriter
}

func NewListAssociationsHandler(lister templateAssociationLister, errWriter errorWriter) ListAssociationsHandler {
	return ListAssociationsHandler{
		lister:      lister,
		errorWriter: errWriter,
	}
}

func (h ListAssociationsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := h.parseTemplateID(req.URL.Path)
	associations, err := h.lister.List(context.Get("database").(DatabaseInterface), templateID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	templateAssociationsDocument := h.mapToJSON(associations)
	writeJSON(w, http.StatusOK, templateAssociationsDocument)
}

func (h ListAssociationsHandler) parseTemplateID(path string) string {
	r := regexp.MustCompile(`\/templates\/(.*)\/associations`)
	matches := r.FindStringSubmatch(path)

	return matches[1]
}

func (h ListAssociationsHandler) mapToJSON(associations []services.TemplateAssociation) map[string][]TemplateAssociation {
	structure := map[string][]TemplateAssociation{
		"associations": {},
	}

	for _, association := range associations {
		structure["associations"] = append(structure["associations"], TemplateAssociation{
			Client:       association.ClientID,
			Notification: association.NotificationID,
		})
	}

	return structure
}
