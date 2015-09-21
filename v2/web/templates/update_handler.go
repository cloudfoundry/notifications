package templates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionSetGetter interface {
	Set(conn collections.ConnectionInterface, template collections.Template) (createdTemplate collections.Template, err error)
	Get(conn collections.ConnectionInterface, templateID, clientID string) (template collections.Template, err error)
}

type UpdateHandler struct {
	templates collectionSetGetter
}

func NewUpdateHandler(templates collectionSetGetter) UpdateHandler {
	return UpdateHandler{
		templates: templates,
	}
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	templateID := splitURL[len(splitURL)-1]

	var updateRequest struct {
		Name     *string          `json:"name"`
		HTML     *string          `json:"html"`
		Text     *string          `json:"text"`
		Subject  *string          `json:"subject"`
		Metadata *json.RawMessage `json:"metadata"`
	}

	err := json.NewDecoder(req.Body).Decode(&updateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "errors": [ "malformed JSON request" ]}`))
		return
	}

	clientID := context.Get("client_id").(string)

	database := context.Get("database").(DatabaseInterface)

	template, err := h.templates.Get(database.Connection(), templateID, clientID)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(fmt.Sprintf(`{ "errors": [ %q ]}`, err)))
		return
	}

	if updateRequest.Name != nil {
		template.Name = *updateRequest.Name
	}

	if updateRequest.HTML != nil {
		template.HTML = *updateRequest.HTML
	}

	if updateRequest.Text != nil {
		template.Text = *updateRequest.Text
	}

	if updateRequest.Subject != nil {
		template.Subject = *updateRequest.Subject
	}

	if updateRequest.Metadata != nil {
		template.Metadata = string(*updateRequest.Metadata)
	}

	if template.Name == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": [ "Template \"name\" field cannot be empty" ] }`))
		return
	}

	if template.Subject == "" {
		template.Subject = "{{.Subject}}"
	}

	if template.HTML == "" && template.Text == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": ["missing either template text or html"] }`))
		return
	}

	template, err = h.templates.Set(database.Connection(), template)
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(fmt.Sprintf(`{ "errors": [ %q ]}`, err)))
		return
	}

	json.NewEncoder(w).Encode(NewTemplateResponse(template))
}
