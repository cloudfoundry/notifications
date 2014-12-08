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

var _ = Describe("SetTemplates", func() {
	var err error
	var handler handlers.SetTemplates
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var context stack.Context
	var updater *fakes.TemplateUpdater
	var errorWriter *fakes.ErrorWriter

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			updater = fakes.NewTemplateUpdater()
			errorWriter = fakes.NewErrorWriter()
			handler = handlers.NewSetTemplates(updater, errorWriter)
			writer = httptest.NewRecorder()
			body := []byte(`{"text": "{{turkey}}", "html": "<p>{{turkey}} gobble</p>"}`)
			request, err = http.NewRequest("PUT", "/deprecated_templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		})

		It("calls set on its setter with appropriate arguments", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(updater.UpdateArgument).To(Equal(models.Template{
				Name: "myTemplateName." + models.UserBodyTemplateName,
				Text: "{{turkey}}",
				HTML: "<p>{{turkey}} gobble</p>",
			}))
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		It("can set a template with an empty text field", func() {
			body := []byte(`{"html": "<p>gobble</p>", "text": ""}`)
			request, err = http.NewRequest("PUT", "/deprecated_templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		It("can set a template with an empty html field", func() {
			body := []byte(`{"html": "", "text": "gobble"}`)
			request, err = http.NewRequest("PUT", "/deprecated_templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("when an errors occurs", func() {
			It("Writes a validation error to the errorwriter when the request is missing the text field", func() {
				body := []byte(`{"html": "<p>gobble</p>"}`)
				request, err = http.NewRequest("PUT", "/deprecated_templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{
					"Request is missing a required field",
				})))
			})

			It("Writes a validation error to the errorwriter when the request is missing the html field", func() {
				body := []byte(`{"text": "gobble"}`)
				request, err = http.NewRequest("PUT", "/deprecated_templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{
					"Request is missing a required field",
				})))
			})

			It("writes a parse error for an invalid request", func() {
				body := []byte(`{"text": forgot to close the curly brace`)
				request, err = http.NewRequest("PUT", "/deprecated_templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ParseError{}))
			})

			It("returns a 500 for all other error cases", func() {
				updater.UpdateError = fmt.Errorf("my new error")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(updater.UpdateError))
			})
		})

		Context("when the template name is malformed", func() {
			It("Writes a validation error when missing a valid ending", func() {
				bad_endings := []string{"user.body", "something_body", "subject.something", "still.missing.something"}

				for _, ending := range bad_endings {
					body := []byte(`{"text": "gobble", "html": "<p>gobble</p>"}`)
					request, err = http.NewRequest("PUT", "/deprecated_templates/"+ending, bytes.NewBuffer(body))
					if err != nil {
						panic(err)
					}
					handler.ServeHTTP(writer, request, context)
					Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{
						fmt.Sprintf("Template has invalid suffix, must end with one of %+v\n", models.TemplateNames),
					})))

				}
			})
		})
	})
})
