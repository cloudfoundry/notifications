package templates

import (
	"encoding/json"
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
		Name     string                 `json:"name"`
		HTML     string                 `json:"html"`
		Text     string                 `json:"text"`
		Subject  string                 `json:"subject"`
		Metadata map[string]interface{} `json:"metadata"`
		ClientID string                 `json:"client_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "errors": [ "invalid json body" ] }`))
		return
	}

	if createRequest.Name == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": ["missing template name"] }`))
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
		createRequest.Metadata = map[string]interface{}{}
	}

	metadata, err := json.Marshal(createRequest.Metadata)
	if err != nil {
		panic(err)
	}

	database := context.Get("database").(DatabaseInterface)

	template, err := h.templates.Set(database.Connection(), collections.Template{
		Name:     createRequest.Name,
		HTML:     createRequest.HTML,
		Text:     createRequest.Text,
		Subject:  createRequest.Subject,
		Metadata: string(metadata),
		ClientID: clientID,
	})

	var decodedMetadata map[string]interface{}
	err = json.Unmarshal([]byte(template.Metadata), &decodedMetadata)
	if err != nil {
		panic(err)
	}

	createResponse, err := json.Marshal(map[string]interface{}{
		"id":       template.ID,
		"name":     template.Name,
		"html":     template.HTML,
		"text":     template.Text,
		"subject":  template.Subject,
		"metadata": decodedMetadata,
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(createResponse)
}
