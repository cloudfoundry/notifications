package templates

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type templateLister interface {
	List(connection collections.ConnectionInterface, clientID string) ([]collections.Template, error)
}

type ListHandler struct {
	templates templateLister
}

func NewListHandler(templates templateLister) ListHandler {
	return ListHandler{
		templates: templates,
	}
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := context.Get("database").(collections.DatabaseInterface).Connection()

	clientID := context.Get("client_id").(string)

	templates, err := h.templates.List(connection, clientID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors": [ %q ] }`, err)
		return
	}

	responseList := []interface{}{}

	for _, template := range templates {
		var metadata map[string]interface{}

		err = json.Unmarshal([]byte(template.Metadata), &metadata)
		if err != nil {
			panic(err)
		}

		responseList = append(responseList, map[string]interface{}{
			"id":       template.ID,
			"name":     template.Name,
			"text":     template.Text,
			"html":     template.HTML,
			"subject":  template.Subject,
			"metadata": metadata,
		})
	}

	listResponse, _ := json.Marshal(map[string]interface{}{
		"templates": responseList,
	})

	w.Write(listResponse)

}
