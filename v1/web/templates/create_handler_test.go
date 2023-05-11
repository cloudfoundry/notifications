package templates_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/collections"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateHandler", func() {
	var (
		err         error
		handler     templates.CreateHandler
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		creator     *mocks.TemplateCreator
		errorWriter *mocks.ErrorWriter
		connection  *mocks.Connection
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			creator = mocks.NewTemplateCreator()
			creator.CreateCall.Returns.Template = collections.Template{
				ID:       "template-guid",
				Name:     "Emergency Template",
				Text:     "Message to: {{.To}}. Raptor Alert.",
				HTML:     "<p>{{.ClientID}} you should run.</p>",
				Subject:  "Raptor Containment Unit Breached",
				Metadata: "{}",
			}

			errorWriter = mocks.NewErrorWriter()
			writer = httptest.NewRecorder()
			body := bytes.NewBuffer([]byte{})
			err := json.NewEncoder(body).Encode(map[string]interface{}{
				"name":    "Emergency Template",
				"text":    "Message to: {{.To}}. Raptor Alert.",
				"html":    "<p>{{.ClientID}} you should run.</p>",
				"subject": "Raptor Containment Unit Breached",
			})
			Expect(err).NotTo(HaveOccurred())

			connection = mocks.NewConnection()
			database := mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = connection

			context = stack.NewContext()
			context.Set("database", database)

			request, err = http.NewRequest("POST", "/templates", body)
			Expect(err).NotTo(HaveOccurred())

			handler = templates.NewCreateHandler(creator, errorWriter)
		})

		It("calls create on its Creator with the correct arguments", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(creator.CreateCall.Receives.Connection).To(Equal(connection))
			Expect(creator.CreateCall.Receives.Template).To(Equal(collections.Template{
				Name:     "Emergency Template",
				Text:     "Message to: {{.To}}. Raptor Alert.",
				HTML:     "<p>{{.ClientID}} you should run.</p>",
				Subject:  "Raptor Containment Unit Breached",
				Metadata: "{}",
			}))

			Expect(writer.Code).To(Equal(http.StatusCreated))
			Expect(writer.Body.String()).To(MatchJSON(`{"template_id":"template-guid"}`))
		})

		Context("when an errors occurs", func() {
			It("Writes a validation error to the errorwriter when the request is missing the name field", func() {
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer([]byte(`{"html": "<p>gobble</p>"}`)))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.ValidationError{Err: valiant.RequiredFieldError{ErrorMessage: "Missing required field 'name'"}}))
			})

			It("Writes a validation error to the errorwriter when the request is missing the html field", func() {
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer([]byte(`{"name": "gobble"}`)))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.ValidationError{Err: valiant.RequiredFieldError{ErrorMessage: "Missing required field 'html'"}}))
			})

			It("writes a parse error for an invalid request", func() {
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer([]byte(`{"name":"foobar", "html": forgot to close the curly brace`)))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ParseError{}))
			})

			It("returns a 500 for all other error cases", func() {
				creator.CreateCall.Returns.Error = errors.New("my new error")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.TemplateCreateError{}))
			})
		})
	})
})
