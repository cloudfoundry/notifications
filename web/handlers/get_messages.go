package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/ryanmoran/stack"
)

type GetMessages struct {
	finder      MessageFinderInterface
	errorWriter ErrorWriterInterface
}

type MessageFinderInterface interface {
	Find(models.DatabaseInterface, string) (services.Message, error)
}

func NewGetMessages(messageFinder MessageFinderInterface, errorWriter ErrorWriterInterface) GetMessages {
	return GetMessages{
		finder:      messageFinder,
		errorWriter: errorWriter,
	}
}

func (handler GetMessages) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	messageID := strings.Split(req.URL.Path, "/messages/")[1]

	message, err := handler.finder.Find(context.Get("database").(models.DatabaseInterface), messageID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	var document struct {
		Status string `json:"status"`
	}
	document.Status = message.Status

	writeJSON(w, http.StatusOK, document)
}
