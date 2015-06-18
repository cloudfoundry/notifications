package handlers

import (
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type ListTemplateAssociations struct {
	lister      services.TemplateAssociationListerInterface
	errorWriter ErrorWriterInterface
}

type TemplateAssociation struct {
	Client       string `json:"client"`
	Notification string `json:"notification,omitempty"`
}

func NewListTemplateAssociations(lister services.TemplateAssociationListerInterface, errorWriter ErrorWriterInterface) ListTemplateAssociations {
	return ListTemplateAssociations{
		lister:      lister,
		errorWriter: errorWriter,
	}
}

func (handler ListTemplateAssociations) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := handler.parseTemplateID(req.URL.Path)
	associations, err := handler.lister.List(context.Get("database").(models.DatabaseInterface), templateID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	templateAssociationsDocument := handler.mapToJSON(associations)
	writeJSON(w, http.StatusOK, templateAssociationsDocument)
}

func (handler ListTemplateAssociations) parseTemplateID(path string) string {
	r := regexp.MustCompile(`\/templates\/(.*)\/associations`)
	matches := r.FindStringSubmatch(path)

	return matches[1]
}

func (handler ListTemplateAssociations) mapToJSON(associations []services.TemplateAssociation) map[string][]TemplateAssociation {
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
