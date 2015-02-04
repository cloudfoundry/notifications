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
	case postal.UAAScopesError:
		writer.write(w, 422, []string{err.Error()})
	case utilities.CCDownError:
		writer.write(w, http.StatusBadGateway, []string{err.Error()})
	case utilities.CCNotFoundError:
		writer.write(w, http.StatusNotFound, []string{err.Error()})
	case postal.UAADownError:
		writer.write(w, http.StatusBadGateway, []string{err.Error()})
	case postal.UAAGenericError:
		writer.write(w, http.StatusBadGateway, []string{err.Error()})
	case postal.TemplateLoadError:
		writer.write(w, http.StatusInternalServerError, []string{err.Error()})
	case params.TemplateCreateError:
		writer.write(w, http.StatusInternalServerError, []string{err.Error()})
	case models.TemplateFindError:
		writer.write(w, http.StatusNotFound, []string{err.Error()})
	case models.TemplateUpdateError:
		writer.write(w, http.StatusInternalServerError, []string{err.Error()})
	case params.ParseError:
		writer.write(w, http.StatusBadRequest, []string{err.Error()})
	case params.ValidationError:
		writer.write(w, 422, err.(params.ValidationError).Errors())
	case params.SchemaError:
		writer.write(w, http.StatusBadRequest, []string{err.Error()})
	case postal.CriticalNotificationError:
		writer.write(w, 422, []string{err.Error()})
	case models.DuplicateRecordError:
		writer.write(w, http.StatusConflict, []string{err.Error()})
	case models.RecordNotFoundError:
		writer.write(w, http.StatusNotFound, []string{err.Error()})
	case models.TransactionCommitError:
		writer.write(w, http.StatusInternalServerError, []string{err.Error()})
	case strategies.DefaultScopeError:
		writer.write(w, http.StatusNotAcceptable, []string{err.Error()})
	case services.TemplateAssignmentError:
		writer.write(w, 422, []string{err.Error()})
	case MissingUserTokenError:
		writer.write(w, 422, []string{err.Error()})
	default:
		panic(err)
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
