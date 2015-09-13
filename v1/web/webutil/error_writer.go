package webutil

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type ErrorWriterInterface interface {
	Write(http.ResponseWriter, error)
}

type ErrorWriter struct{}

func NewErrorWriter() ErrorWriter {
	return ErrorWriter{}
}

func (writer ErrorWriter) Write(w http.ResponseWriter, err error) {
	switch err.(type) {
	case UAAScopesError, postal.CriticalNotificationError, services.TemplateAssignmentError, MissingUserTokenError, ValidationError:
		writer.write(w, 422, []string{err.Error()})
	case services.CCDownError, postal.UAADownError, postal.UAAGenericError:
		writer.write(w, http.StatusBadGateway, []string{err.Error()})
	case services.CCNotFoundError, models.TemplateFindError, models.NotFoundError:
		writer.write(w, http.StatusNotFound, []string{err.Error()})
	case postal.TemplateLoadError, TemplateCreateError, models.TemplateUpdateError, models.TransactionCommitError:
		writer.write(w, http.StatusInternalServerError, []string{err.Error()})
	case ParseError, SchemaError:
		writer.write(w, http.StatusBadRequest, []string{err.Error()})
	case models.DuplicateError:
		writer.write(w, http.StatusConflict, []string{err.Error()})
	case services.DefaultScopeError:
		writer.write(w, http.StatusNotAcceptable, []string{err.Error()})
	default:
		panic(err) // This panic will trigger the Stack recovery handler
	}
}

func (writer ErrorWriter) write(w http.ResponseWriter, code int, errors []string) {
	response, err := json.Marshal(map[string][]string{
		"errors": errors,
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(code)
	w.Write(response)
}
