package handlers_test

import (
	"bytes"
	"fmt"
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

var _ = Describe("CreateTemplates", func() {
	var err error
	var handler handlers.CreateTemplates
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var context stack.Context
	var creator *fakes.TemplateCreator
	var errorWriter *fakes.ErrorWriter

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			creator = fakes.NewTemplateCreator()
			errorWriter = fakes.NewErrorWriter()
			handler = handlers.NewCreateTemplates(creator, errorWriter)
			writer = httptest.NewRecorder()
			body := []byte(`{"name": "Emergency Template", "text": "Message to: {{.To}}. Raptor Alert.", "html": "<p>{{.ClientID}} you should run.</p>", "subject": "Raptor Containment Unit Breached"}`)
			request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		})

		It("calls create on its Creator with the correct arguements", func() {
			handler.ServeHTTP(writer, request, context)
			body := string(writer.Body.Bytes())

			Expect(creator.CreateArgument).To(Equal(models.Template{
				Name:    "Emergency Template",
				Text:    "Message to: {{.To}}. Raptor Alert.",
				HTML:    "<p>{{.ClientID}} you should run.</p>",
				Subject: "Raptor Containment Unit Breached",
			}))
			Expect(writer.Code).To(Equal(http.StatusCreated))
			Expect(body).To(Equal(`{"template-id":"guid"}`))
		})

		Context("when an errors occurs", func() {
			It("Writes a validation error to the errorwriter when the request is missing the name field", func() {
				body := []byte(`{"html": "<p>gobble</p>"}`)
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{
					"Request is missing the required field: name",
				})))
			})

			It("Writes a validation error to the errorwriter when the request is missing the html field", func() {
				body := []byte(`{"name": "gobble"}`)
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{
					"Request is missing the required field: html",
				})))
			})

			It("writes a parse error for an invalid request", func() {
				body := []byte(`{"name":"foobar", "html": forgot to close the curly brace`)
				request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ParseError{}))
			})

			It("returns a 500 for all other error cases", func() {
				creator.CreateError = fmt.Errorf("my new error")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.TemplateCreateError{}))
			})
		})
	})
})
