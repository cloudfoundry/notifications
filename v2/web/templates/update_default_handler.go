package templates

import (
	"encoding/json"
	"fmt"
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
		Name     *string          `json:"name"`
		HTML     *string          `json:"html"`
		Text     *string          `json:"text"`
		Subject  *string          `json:"subject"`
		Metadata *json.RawMessage `json:"metadata"`
	}

	err := json.NewDecoder(req.Body).Decode(&updateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "errors": [ "malformed JSON request" ] }`))
		return
	}

	database := context.Get("database").(DatabaseInterface)

	template, err := h.templatesCollection.Get(database.Connection(), "default", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "errors": [%q] }`, err)
		return
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
		template.Metadata = string(*updateRequest.Metadata)
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "errors": [%q] }`, err)
		return
	}

	json.NewEncoder(w).Encode(NewTemplateResponse(template))
}
