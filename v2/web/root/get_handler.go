package root

import (
	"encoding/json"
	"net/http"

	"github.com/ryanmoran/stack"
)

type GetHandler struct {
}

func NewGetHandler() GetHandler {
	return GetHandler{}
}

func (handler GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	output, err := json.Marshal(map[string]interface{}{
		"_links": map[string]interface{}{
			"self": map[string]string{
				"href": "/",
			},
			"senders": map[string]string{
				"href": "/senders",
			},
			"templates": map[string]string{
				"href": "/templates",
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(output)
}
