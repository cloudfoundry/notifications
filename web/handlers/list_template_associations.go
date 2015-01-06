package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

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
	associations, err := handler.lister.List(templateID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	response, err := json.Marshal(handler.mapToJSON(associations))
	if err != nil {
		panic(err)
	}

	w.Write(response)
}

func (handler ListTemplateAssociations) parseTemplateID(path string) string {
	r := regexp.MustCompile(`\/templates\/(.*)\/associations`)
	matches := r.FindStringSubmatch(path)

	return matches[1]
}

func (handler ListTemplateAssociations) mapToJSON(associations []services.TemplateAssociation) map[string][]TemplateAssociation {
	structure := map[string][]TemplateAssociation{
		"associations": []TemplateAssociation{},
	}

	for _, association := range associations {
		structure["associations"] = append(structure["associations"], TemplateAssociation{
			Client:       association.ClientID,
			Notification: association.NotificationID,
		})
	}

	return structure
}
