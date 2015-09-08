package templates

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type templateDeleter interface {
	Delete(database services.DatabaseInterface, templateID string) error
}

type DeleteHandler struct {
	deleter     templateDeleter
	errorWriter errorWriter
}

func NewDeleteHandler(deleter templateDeleter, errWriter errorWriter) DeleteHandler {
	return DeleteHandler{
		deleter:     deleter,
		errorWriter: errWriter,
	}
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	err := h.deleter.Delete(context.Get("database").(DatabaseInterface), templateID)
	if err != nil {
		h.errorWriter.Write(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
