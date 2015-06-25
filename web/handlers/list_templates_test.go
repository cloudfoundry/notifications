package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListTemplates", func() {
	var (
		handler       handlers.ListTemplates
		request       *http.Request
		writer        *httptest.ResponseRecorder
		context       stack.Context
		lister        *fakes.TemplateLister
		errorWriter   *fakes.ErrorWriter
		testTemplates map[string]services.TemplateSummary
		database      *fakes.Database
	)

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

			database = fakes.NewDatabase()
			context = stack.NewContext()
			context.Set("database", database)
		})

		Context("when the lister does not error", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest("GET", "/templates", bytes.NewBuffer([]byte{}))
				Expect(err).NotTo(HaveOccurred())
			})

			It("calls list on its lister", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(lister.ListCall.Arguments).To(ConsistOf([]interface{}{database}))
			})

			It("writes out the lister's response", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(http.StatusOK))

				var templates map[string]services.TemplateSummary

				err := json.Unmarshal(writer.Body.Bytes(), &templates)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(templates)).To(Equal(4))
				Expect(templates).To(Equal(testTemplates))
			})

		})

		Context("when the lister errors", func() {
			BeforeEach(func() {
				lister.ListCall.Error = errors.New("BOOM!")

				var err error
				request, err = http.NewRequest("GET", "/templates", bytes.NewBuffer([]byte{}))
				Expect(err).NotTo(HaveOccurred())
			})

			It("writes the error to the errorWriter", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))
			})
		})

	})
})
