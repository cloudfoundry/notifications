package templates_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateHandler", func() {
	var (
		handler             templates.UpdateHandler
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

		context = stack.NewContext()
		context.Set("database", database)
		context.Set("client_id", "some-client-id")

		writer = httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]interface{}{
			"name":     "an interesting template",
			"text":     "template text",
			"html":     "template html",
			"subject":  "template subject",
			"metadata": map[string]string{"template": "metadata"},
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("PUT", "/templates/some-template-id", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		templatesCollection = mocks.NewTemplatesCollection()
		templatesCollection.GetCall.Returns.Template = collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			HTML:     "template html",
			Text:     "template text",
			Subject:  "template subject",
			Metadata: `{"template": "metadata"}`,
			ClientID: "some-client-id",
		}
		templatesCollection.SetCall.Returns.Template = collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			HTML:     "template html",
			Text:     "template text",
			Subject:  "template subject",
			Metadata: `{"template": "metadata"}`,
		}

		handler = templates.NewUpdateHandler(templatesCollection)
	})

	It("updates a template", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-template-id",
			"name": "an interesting template",
			"text": "template text",
			"html": "template html",
			"subject": "template subject",
			"metadata": {
				"template": "metadata"
			}
		}`))

		Expect(templatesCollection.SetCall.Receives.Connection).To(Equal(conn))
		Expect(templatesCollection.SetCall.Receives.Template).To(Equal(collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			HTML:     "template html",
			Text:     "template text",
			Subject:  "template subject",
			Metadata: `{"template":"metadata"}`,
			ClientID: "some-client-id",
		}))
	})

	Context("when updating the default template", func() {
		It("saves an empty client ID", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"name":     "a default template",
				"text":     "new text",
				"html":     "default html",
				"subject":  "default subject",
				"metadata": map[string]string{"template": "default"},
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/default", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			templatesCollection.GetCall.Returns.Template = collections.Template{
				ID:       "default",
				Name:     "a default template",
				HTML:     "default html",
				Text:     "default text",
				Subject:  "default subject",
				Metadata: `{"template": "default"}`,
			}

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))

			Expect(templatesCollection.SetCall.Receives.Connection).To(Equal(conn))
			Expect(templatesCollection.SetCall.Receives.Template).To(Equal(collections.Template{
				ID:       "default",
				Name:     "a default template",
				HTML:     "default html",
				Text:     "new text",
				Subject:  "default subject",
				Metadata: `{"template":"default"}`,
				ClientID: "",
			}))
		})
	})

	Context("when omitting fields", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/some-template-id", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			templatesCollection.SetCall.Returns.Template = collections.Template{
				ID:       "some-template-id",
				Name:     "an interesting template",
				HTML:     "template html",
				Text:     "template text",
				Subject:  "template subject",
				Metadata: `{"template": "metadata"}`,
			}
		})

		It("does not change any of the existing fields", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"id": "some-template-id",
				"name": "an interesting template",
				"text": "template text",
				"html": "template html",
				"subject": "template subject",
				"metadata": {
					"template": "metadata"
				}
			}`))

			Expect(templatesCollection.SetCall.Receives.Connection).To(Equal(conn))
			Expect(templatesCollection.SetCall.Receives.Template).To(Equal(collections.Template{
				ID:       "some-template-id",
				Name:     "an interesting template",
				HTML:     "template html",
				Text:     "template text",
				Subject:  "template subject",
				Metadata: `{"template": "metadata"}`,
				ClientID: "some-client-id",
			}))
		})
	})

	Context("when clearing the subject field", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"name":     "an interesting template",
				"text":     "template text",
				"html":     "template html",
				"subject":  "",
				"metadata": map[string]string{"template": "metadata"},
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/some-template-id", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			templatesCollection.SetCall.Returns.Template = collections.Template{
				ID:       "some-template-id",
				Name:     "an interesting template",
				HTML:     "template html",
				Text:     "template text",
				Subject:  "{{.Subject}}",
				Metadata: `{"template": "metadata"}`,
			}
		})

		It("repopulates the default subject field", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"id": "some-template-id",
				"name": "an interesting template",
				"text": "template text",
				"html": "template html",
				"subject": "{{.Subject}}",
				"metadata": {
					"template": "metadata"
				}
			}`))

			Expect(templatesCollection.SetCall.Receives.Connection).To(Equal(conn))
			Expect(templatesCollection.SetCall.Receives.Template).To(Equal(collections.Template{
				ID:       "some-template-id",
				Name:     "an interesting template",
				HTML:     "template html",
				Text:     "template text",
				Subject:  "{{.Subject}}",
				Metadata: `{"template":"metadata"}`,
				ClientID: "some-client-id",
			}))
		})
	})

	Context("when the template does not exist", func() {
		It("returns a 404 and and error message", func() {
			templatesCollection.GetCall.Returns.Error = collections.NotFoundError{errors.New("not found")}

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["not found"]
			}`))
		})
	})

	Context("when the name field is empty", func() {
		It("returns a 422 with an error message", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"name":     "",
				"text":     "template text",
				"html":     "template html",
				"subject":  "template subject",
				"metadata": map[string]string{"template": "metadata"},
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/some-template-id", bytes.NewBuffer(requestBody))
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
			templatesCollection.GetCall.Returns.Template = collections.Template{
				ID:       "some-template-id",
				Name:     "an interesting template",
				HTML:     "Something",
				Text:     "",
				Subject:  "template subject",
				Metadata: `{"template": "metadata"}`,
			}

			requestBody, err := json.Marshal(map[string]interface{}{
				"html": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", "/templates/some-template-id", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing either template text or html"]
			}`))
		})
	})

	Context("failure cases", func() {
		It("returns a 400 with an error message if the request JSON is malformed", func() {
			var err error
			request, err = http.NewRequest("PUT", "/templates/some-template-id", bytes.NewBuffer([]byte("%%%")))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["malformed JSON request"]
			}`))
		})

		It("returns a 500 with an error message if getting the template fails", func() {
			templatesCollection.GetCall.Returns.Error = errors.New("something bad happened")

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["something bad happened"]
			}`))
		})

		It("returns a 500 with an error message if setting the template fails", func() {
			templatesCollection.SetCall.Returns.Error = errors.New("failed to talk to the db")

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["failed to talk to the db"]
			}`))
		})
	})
})
