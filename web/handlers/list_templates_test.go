package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListTemplates", func() {
	var handler handlers.ListTemplates
	var request *http.Request
	var writer *httptest.ResponseRecorder
	var context stack.Context
	var lister *fakes.TemplateLister
	var errorWriter *fakes.ErrorWriter
	var testTemplates map[string]services.TemplateSummary

	Describe("ServeHTTP", func() {
		BeforeEach(func() {

			testTemplates = map[string]services.TemplateSummary{
				"chewbaca-guid": {
					Name: "Star Wars",
				},
				"giant-friendly-robot-guid": {
					Name: "Big Hero 6",
				},
				"boring-template-guid": {
					Name: "Blah",
				},
				"starvation-guid": {
					Name: "Hungry Play",
				},
			}

			lister = fakes.NewTemplateLister()
			lister.Templates = testTemplates

			writer = httptest.NewRecorder()
			errorWriter = fakes.NewErrorWriter()
			handler = handlers.NewListTemplates(lister, errorWriter)
		})

		Context("When the lister does not error", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest("GET", "/templates", bytes.NewBuffer([]byte{}))
				if err != nil {
					panic(err)
				}
			})

			It("Calls list on its lister", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(lister.ListWasCalled).To(BeTrue())
			})

			It("writes out the lister's response", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(http.StatusOK))

				var templates map[string]services.TemplateSummary

				err := json.Unmarshal(writer.Body.Bytes(), &templates)
				if err != nil {
					panic(err)
				}

				Expect(len(templates)).To(Equal(4))
				Expect(templates).To(Equal(testTemplates))
			})

		})

		Context("When the lister errors", func() {
			BeforeEach(func() {
				lister.ListError = errors.New("BOOM!")

				var err error
				request, err = http.NewRequest("GET", "/templates", bytes.NewBuffer([]byte{}))
				if err != nil {
					panic(err)
				}
			})

			It("writes the error to the errorWriter", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))
			})
		})

	})
})
