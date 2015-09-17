package templates_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateDefaultHandler", func() {
	var (
		handler             templates.UpdateDefaultHandler
		context             stack.Context
		writer              *httptest.ResponseRecorder
		request             *http.Request
		templatesCollection *mocks.TemplatesCollection
		conn                *mocks.Connection
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		database := mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		templatesCollection = mocks.NewTemplatesCollection()
		templatesCollection.GetCall.Returns.Template = collections.Template{
			ID:       "default",
			Name:     "a default template",
			HTML:     "default html",
			Text:     "default text",
			Subject:  "default subject",
			Metadata: `{"template": "default"}`,
		}

		context = stack.NewContext()
		context.Set("database", database)
		context.Set("client_id", "admin-client")

		writer = httptest.NewRecorder()

		handler = templates.NewUpdateDefaultHandler(templatesCollection)
	})

	It("updates the default template", func() {
		requestBody, err := json.Marshal(map[string]interface{}{
			"name":     "new template name",
			"html":     "new html",
			"text":     "new text",
			"subject":  "new subject",
			"metadata": map[string]string{"template": "new"},
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("PUT", "/templates/default", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		templatesCollection.SetCall.Returns.Template = collections.Template{
			ID:       "default",
			Name:     "new template name",
			HTML:     "new html",
			Text:     "new text",
			Subject:  "new subject",
			Metadata: `{"template": "new"}`,
		}

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body).To(MatchJSON(`{
			"id":       "default",
			"name":     "new template name",
			"html":     "new html",
			"text":     "new text",
			"subject":  "new subject",
			"metadata": {"template":"new"}
		}`))

		Expect(templatesCollection.SetCall.Receives.Connection).To(Equal(conn))
		Expect(templatesCollection.SetCall.Receives.Template).To(Equal(collections.Template{
			ID:       "default",
			Name:     "new template name",
			HTML:     "new html",
			Text:     "new text",
			Subject:  "new subject",
			Metadata: `{"template":"new"}`,
			ClientID: "",
		}))
	})

	Context("when omitting fields", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/default", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			templatesCollection.SetCall.Returns.Template = collections.Template{
				ID:       "default",
				Name:     "a default template",
				HTML:     "default html",
				Text:     "default text",
				Subject:  "default subject",
				Metadata: `{"template": "default"}`,
			}
		})

		It("does not change any of the existing fields", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"id":       "default",
				"name":     "a default template",
				"html":     "default html",
				"text":     "default text",
				"subject":  "default subject",
				"metadata": {"template":"default"}
			}`))

			Expect(templatesCollection.SetCall.Receives.Connection).To(Equal(conn))
			Expect(templatesCollection.SetCall.Receives.Template).To(Equal(collections.Template{
				ID:       "default",
				Name:     "a default template",
				HTML:     "default html",
				Text:     "default text",
				Subject:  "default subject",
				Metadata: `{"template": "default"}`,
			}))
		})
	})

	Context("when the subject field is empty", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"subject": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/default", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			templatesCollection.SetCall.Returns.Template = collections.Template{
				ID:       "default",
				Name:     "a default template",
				HTML:     "default html",
				Text:     "default text",
				Subject:  "{{.Subject}}",
				Metadata: `{"template": "default"}`,
			}
		})

		It("repopulates the default subject field", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"id":       "default",
				"name":     "a default template",
				"html":     "default html",
				"text":     "default text",
				"subject":  "{{.Subject}}",
				"metadata": {"template":"default"}
			}`))

			Expect(templatesCollection.SetCall.Receives.Connection).To(Equal(conn))
			Expect(templatesCollection.SetCall.Receives.Template).To(Equal(collections.Template{
				ID:       "default",
				Name:     "a default template",
				HTML:     "default html",
				Text:     "default text",
				Subject:  "{{.Subject}}",
				Metadata: `{"template": "default"}`,
			}))
		})
	})

	Context("when the name field is empty", func() {
		It("returns a 422 with an error message", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"name": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/default", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["Template \"name\" field cannot be empty"]
			}`))
		})
	})

	Context("when the html and text field would be empty", func() {
		It("returns a 422 with an error message", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"text": "",
				"html": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/default", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing either template text or html"]
			}`))
		})
	})
})
