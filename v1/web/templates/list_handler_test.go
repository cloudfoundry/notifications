package templates_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListHandler", func() {
	var (
		handler       templates.ListHandler
		request       *http.Request
		writer        *httptest.ResponseRecorder
		context       stack.Context
		lister        *mocks.TemplateLister
		errorWriter   *mocks.ErrorWriter
		testTemplates map[string]services.TemplateSummary
		database      *mocks.Database
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

			lister = mocks.NewTemplateLister()
			lister.ListCall.Returns.TemplateSummaries = testTemplates

			writer = httptest.NewRecorder()
			errorWriter = mocks.NewErrorWriter()

			database = mocks.NewDatabase()
			context = stack.NewContext()
			context.Set("database", database)

			handler = templates.NewListHandler(lister, errorWriter)
		})

		Context("when the lister does not error", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest("GET", "/templates", bytes.NewBuffer([]byte{}))
				Expect(err).NotTo(HaveOccurred())
			})

			It("calls list on its lister", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(lister.ListCall.Receives.Database).To(Equal(database))
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
				lister.ListCall.Returns.Error = errors.New("BOOM!")

				var err error
				request, err = http.NewRequest("GET", "/templates", bytes.NewBuffer([]byte{}))
				Expect(err).NotTo(HaveOccurred())
			})

			It("writes the error to the errorWriter", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BOOM!")))
			})
		})

	})
})
