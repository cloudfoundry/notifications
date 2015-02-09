package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type GetMessages struct {
	Finder      MessageFinderInterface
	errorWriter ErrorWriterInterface
}

type MessageFinderInterface interface {
	Find(string) (services.Message, error)
}

func NewGetMessages(messageFinder MessageFinderInterface, errorWriter ErrorWriterInterface) GetMessages {
	return GetMessages{
		Finder:      messageFinder,
		errorWriter: errorWriter,
	}
}

func (handler GetMessages) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.messages.get",
	}).Log()

	messageID := strings.Split(req.URL.Path, "/messages/")[1]

	message, err := handler.Finder.Find(messageID)
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
