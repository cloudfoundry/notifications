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
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetHandler", func() {
	var (
		handler     templates.GetHandler
		request     *http.Request
		writer      *httptest.ResponseRecorder
		context     stack.Context
		finder      *mocks.TemplateFinder
		errorWriter *mocks.ErrorWriter
		database    *mocks.Database
	)

	Describe("ServeHTTP", func() {
		var templateID string

		BeforeEach(func() {
			finder = mocks.NewTemplateFinder()
			templateID = "theTemplateID"

			finder.FindByIDCall.Returns.Template = models.Template{
				Name:     "The Name of The Template",
				Subject:  "All about the {{.Subject}}",
				Text:     "the template {{variable}}",
				HTML:     "<p> the template {{variable}} </p>",
				Metadata: `{"hello": "world"}`,
			}
			writer = httptest.NewRecorder()
			errorWriter = mocks.NewErrorWriter()
			database = mocks.NewDatabase()
			context = stack.NewContext()
			context.Set("database", database)

			handler = templates.NewGetHandler(finder, errorWriter)
		})

		Context("when the finder does not error", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest("GET", "/templates/"+templateID, bytes.NewBuffer([]byte{}))
				Expect(err).NotTo(HaveOccurred())
			})

			It("calls find on its finder with appropriate arguments", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(finder.FindByIDCall.Receives.Database).To(Equal(database))
				Expect(finder.FindByIDCall.Receives.TemplateID).To(Equal(templateID))
			})

			It("writes out the finder's response", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(http.StatusOK))

				var template map[string]interface{}

				err := json.Unmarshal(writer.Body.Bytes(), &template)
				if err != nil {
					panic(err)
				}

				Expect(template).To(HaveLen(5))
				Expect(template["name"]).To(Equal("The Name of The Template"))
				Expect(template["subject"]).To(Equal("All about the {{.Subject}}"))
				Expect(template["text"]).To(Equal("the template {{variable}}"))
				Expect(template["html"]).To(Equal("<p> the template {{variable}} </p>"))
				Expect(template["metadata"]).To(Equal(map[string]interface{}{"hello": "world"}))
			})
		})

		Context("when the finder errors", func() {
			BeforeEach(func() {
				finder.FindByIDCall.Returns.Error = errors.New("BOOM!")

				var err error
				request, err = http.NewRequest("GET", "/templates/someTemplateID", bytes.NewBuffer([]byte{}))
				if err != nil {
					panic(err)
				}
			})

			It("writes the error to the errorWriter", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BOOM!")))
			})
		})
	})
})
