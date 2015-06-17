package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
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
	case postal.UAAScopesError, postal.CriticalNotificationError, services.TemplateAssignmentError, MissingUserTokenError:
		writer.write(w, 422, []string{err.Error()})
	case params.ValidationError:
		writer.write(w, 422, err.(params.ValidationError).Errors())
	case utilities.CCDownError, postal.UAADownError, postal.UAAGenericError:
		writer.write(w, http.StatusBadGateway, []string{err.Error()})
	case utilities.CCNotFoundError, models.TemplateFindError, models.RecordNotFoundError:
		writer.write(w, http.StatusNotFound, []string{err.Error()})
	case postal.TemplateLoadError, params.TemplateCreateError, models.TemplateUpdateError, models.TransactionCommitError:
		writer.write(w, http.StatusInternalServerError, []string{err.Error()})
	case params.ParseError, params.SchemaError:
		writer.write(w, http.StatusBadRequest, []string{err.Error()})
	case models.DuplicateRecordError:
		writer.write(w, http.StatusConflict, []string{err.Error()})
	case strategies.DefaultScopeError:
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
