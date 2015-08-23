package templates_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
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
		database    *mocks.Database
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			creator = mocks.NewTemplateCreator()
			creator.CreateCall.Returns.TemplateGUID = "template-guid"
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

			database = mocks.NewDatabase()
			context = stack.NewContext()
			context.Set("database", database)

			request, err = http.NewRequest("POST", "/templates", body)
			Expect(err).NotTo(HaveOccurred())

			handler = templates.NewCreateHandler(creator, errorWriter)
		})

		It("calls create on its Creator with the correct arguments", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(creator.CreateCall.Receives.Database).To(Equal(database))
			Expect(creator.CreateCall.Receives.Template).To(Equal(models.Template{
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
				body := []byte(`{"html": "<p>gobble</p>"}`)
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError([]string{})))
			})

			It("Writes a validation error to the errorwriter when the request is missing the html field", func() {
				body := []byte(`{"name": "gobble"}`)
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError([]string{})))
			})

			It("writes a parse error for an invalid request", func() {
				body := []byte(`{"name":"foobar", "html": forgot to close the curly brace`)
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
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
