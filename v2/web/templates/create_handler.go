package templates

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type collectionSetter interface {
	Set(conn collections.ConnectionInterface, template collections.Template) (createdTemplate collections.Template, err error)
}

type CreateHandler struct {
	templates collectionSetter
}

func NewCreateHandler(templates collectionSetter) CreateHandler {
	return CreateHandler{
		templates: templates,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	var createRequest struct {
		Name     string           `json:"name"`
		HTML     string           `json:"html"`
		Text     string           `json:"text"`
		Subject  string           `json:"subject"`
		Metadata *json.RawMessage `json:"metadata"`
		ClientID string           `json:"client_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "errors": [ "invalid json body" ] }`))
		return
	}

	if createRequest.Name == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": ["Template \"name\" field cannot be empty"] }`))
		return
	}

	if createRequest.HTML == "" && createRequest.Text == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": ["missing either template text or html"] }`))
		return
	}

	clientID := context.Get("client_id").(string)
	if clientID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "errors": ["missing client id"] }`))
		return
	}

	if createRequest.Subject == "" {
		createRequest.Subject = "{{.Subject}}"
	}

	if createRequest.Metadata == nil {
		metadata := json.RawMessage("{}")
		createRequest.Metadata = &metadata
	}

	database := context.Get("database").(DatabaseInterface)

	template, err := h.templates.Set(database.Connection(), collections.Template{
		Name:     createRequest.Name,
		HTML:     createRequest.HTML,
		Text:     createRequest.Text,
		Subject:  createRequest.Subject,
		Metadata: string(*createRequest.Metadata),
		ClientID: clientID,
	})
	if err != nil {
		switch err.(type) {
		case collections.DuplicateRecordError:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, `{"errors": [ %q ]}`, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(NewTemplateResponse(template))
}
