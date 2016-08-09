package groups

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type listHandler struct {
	groups *domain.Groups
}

func (h listHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	list := domain.GroupsList(h.groups.All())

	response, err := json.Marshal(list.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
