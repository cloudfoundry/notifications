package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateNotifications", func() {
	var err error
	var handler handlers.UpdateNotifications
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var context stack.Context
	var updater *fakes.NotificationUpdater
	var errorWriter *fakes.ErrorWriter

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			updater = &fakes.NotificationUpdater{}
			errorWriter = fakes.NewErrorWriter()
			handler = handlers.NewUpdateNotifications(updater, errorWriter)
			writer = httptest.NewRecorder()
			body := []byte(`{"description": "test kind", "critical": false, "template": "template-name"}`)
			request, err = http.NewRequest("PUT", "/clients/this-client/notifications/this-kind", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		})

		It("calls update on its updater with appropriate arguments", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(updater.Notification).To(Equal(models.Kind{
				Description: "test kind",
				Critical:    false,
				TemplateID:  "template-name",
				ClientID:    "this-client",
				ID:          "this-kind",
			}))
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("when an error occurs", func() {
			It("propagates the error returned from the updater into the error writer", func() {
				updater.Error = errors.New("error occurred while updating notification")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(MatchError(errors.New("error occurred while updating notification")))
			})

			It("writes a params validation error when the request is semantically invalid", func() {
				body := []byte(`{"description": "test kind", "template": "template-name"}`)
				request, err = http.NewRequest("PUT", "/clients/this-client/notifications/this-kind", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ValidationError{}))
			})
		})
	})
})
