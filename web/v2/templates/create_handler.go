package templates

import (
	"encoding/json"
	"net/http"

	"github.com/ryanmoran/stack"
)

type CreateHandler struct {
}

func NewCreateHandler(collection interface{}) CreateHandler {
	return CreateHandler{}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	createRequest := map[string]interface{}{}
	err := json.NewDecoder(req.Body).Decode(&createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "errors": [ "invalid json body" ] }`))
		return
	}

	if createRequest["name"] == nil || createRequest["name"] == "" {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": ["missing template name"] }`))
		return
	}

	if (createRequest["html"] == nil || createRequest["html"] == "") && (createRequest["text"] == nil || createRequest["text"] == "") {
		w.WriteHeader(422)
		w.Write([]byte(`{ "errors": ["missing either template text or html"] }`))
		return
	}

	if context.Get("client_id").(string) == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "errors": ["missing client id"] }`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	createRequest["id"] = "some-template-id"
	if createRequest["html"] == nil {
		createRequest["html"] = ""
	}
	if createRequest["text"] == nil {
		createRequest["text"] = ""
	}
	if createRequest["subject"] == nil || createRequest["subject"] == "" {
		createRequest["subject"] = "{{.Subject}}"
	}
	if createRequest["metadata"] == nil {
		createRequest["metadata"] = map[string]interface{}{}
	}

	createResponse, err := json.Marshal(createRequest)
	if err != nil {
		panic(err)
	}
	w.Write(createResponse)
}
