package templates

import (
	"encoding/json"
	"net/http"

	"github.com/ryanmoran/stack"
)

type UpdateDefaultHandler struct {
	templatesCollection collectionSetGetter
}

func NewUpdateDefaultHandler(collection collectionSetGetter) UpdateDefaultHandler {
	return UpdateDefaultHandler{
		templatesCollection: collection,
	}
}

func (h UpdateDefaultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	var updateRequest struct {
		Name     *string                `json:"name"`
		HTML     *string                `json:"html"`
		Text     *string                `json:"text"`
		Subject  *string                `json:"subject"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	err := json.NewDecoder(req.Body).Decode(&updateRequest)
	if err != nil {
		panic(err)
	}

	database := context.Get("database").(DatabaseInterface)

	template, err := h.templatesCollection.Get(database.Connection(), "default", "")
	if err != nil {
		panic(err)
	}

	if updateRequest.Name != nil {
		template.Name = *updateRequest.Name
	}

	if updateRequest.Text != nil {
		template.Text = *updateRequest.Text
	}

	if updateRequest.HTML != nil {
		template.HTML = *updateRequest.HTML
	}

	if updateRequest.Subject != nil {
		template.Subject = *updateRequest.Subject
	}

	if updateRequest.Metadata != nil {
		metadata, err := json.Marshal(updateRequest.Metadata)
		if err != nil {
			panic(err)
		}

		template.Metadata = string(metadata)
	}

	if template.Name == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": [ "Template \"name\" field cannot be empty" ] }`))
		return
	}

	if template.Text == "" && template.HTML == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": [ "missing either template text or html" ] }`))
		return
	}

	if template.Subject == "" {
		template.Subject = "{{.Subject}}"
	}

	template, err = h.templatesCollection.Set(database.Connection(), template)
	if err != nil {
		panic(err)
	}

	var decodedMetadata map[string]interface{}
	err = json.Unmarshal([]byte(template.Metadata), &decodedMetadata)
	if err != nil {
		panic(err)
	}

	updateResponse, err := json.Marshal(map[string]interface{}{
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

	w.Write(updateResponse)
}
