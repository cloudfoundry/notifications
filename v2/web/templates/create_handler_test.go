package templates_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateHandler", func() {
	var (
		handler             templates.CreateHandler
		templatesCollection *fakes.TemplatesCollection
		context             stack.Context
		writer              *httptest.ResponseRecorder
		request             *http.Request
		database            *fakes.Database
	)

	BeforeEach(func() {
		context = stack.NewContext()
		context.Set("client_id", "some-client-id")

		database = fakes.NewDatabase()
		context.Set("database", database)

		writer = httptest.NewRecorder()

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":     "an interesting template",
			"text":     "template text",
			"html":     "template html",
			"subject":  "template subject",
			"metadata": map[string]string{"template": "metadata"},
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		templatesCollection = fakes.NewTemplatesCollection()
		handler = templates.NewCreateHandler(templatesCollection)
	})

	It("creates a template", func() {
		templatesCollection.SetCall.Returns.Template = collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			HTML:     "template html",
			Text:     "template text",
			Subject:  "template subject",
			Metadata: `{"template": "metadata"}`,
		}

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusCreated))
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
	})

	It("creates a template with only name and text", func() {
		templatesCollection.SetCall.Returns.Template = collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			Text:     "this is my text",
			Subject:  "{{.Subject}}",
			Metadata: "{}",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name": "an interesting template",
			"text": "this is my text",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(requestBody))

		handler.ServeHTTP(writer, request, context)

		Expect(templatesCollection.SetCall.Receives.Template.Subject).To(Equal("{{.Subject}}"))
		Expect(templatesCollection.SetCall.Receives.Template.Metadata).To(Equal("{}"))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-template-id",
			"name": "an interesting template",
			"text": "this is my text",
			"html": "",
			"subject": "{{.Subject}}",
			"metadata": {}
		}`))
	})

	It("creates a template with only name and html", func() {
		templatesCollection.SetCall.Returns.Template = collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			HTML:     "template html",
			Subject:  "{{.Subject}}",
			Metadata: "{}",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name": "an interesting template",
			"html": "template html",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(requestBody))

		handler.ServeHTTP(writer, request, context)

		Expect(templatesCollection.SetCall.Receives.Template.Subject).To(Equal("{{.Subject}}"))
		Expect(templatesCollection.SetCall.Receives.Template.Metadata).To(Equal("{}"))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-template-id",
			"name": "an interesting template",
			"text": "",
			"html": "template html",
			"subject": "{{.Subject}}",
			"metadata": {}
		}`))
	})

	It("defaults subject when it is empty string", func() {
		templatesCollection.SetCall.Returns.Template = collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			HTML:     "template html",
			Subject:  "{{.Subject}}",
			Metadata: "{}",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":    "an interesting template",
			"html":    "template html",
			"subject": "",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/templates", bytes.NewBuffer(requestBody))

		handler.ServeHTTP(writer, request, context)

		Expect(templatesCollection.SetCall.Receives.Template.Subject).To(Equal("{{.Subject}}"))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-template-id",
			"name": "an interesting template",
			"text": "",
			"html": "template html",
			"subject": "{{.Subject}}",
			"metadata": {}
		}`))
	})

	Context("failure cases", func() {
		It("returns a 400 when the JSON cannot be unmarshalled", func() {
			var err error
			request, err = http.NewRequest("POST", "/templates", strings.NewReader("%%%"))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["invalid json body"]
			}`))
		})

		It("returns a 422 when the request does not include a template name", func() {
			var err error
			request, err = http.NewRequest("POST", "/templates", strings.NewReader("{}"))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing template name"]
			}`))
		})

		It("returns a 422 when the request template name is empty", func() {
			var err error
			request, err = http.NewRequest("POST", "/templates", strings.NewReader(`{"name": ""}`))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing template name"]
			}`))
		})

		It("returns a 422 when the request does not include either a text or html body", func() {
			var err error
			request, err = http.NewRequest("POST", "/templates", strings.NewReader(`{
				"name": "a cool template"
			}`))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing either template text or html"]
			}`))
		})

		It("returns a 422 when the request includes an empty text body and no html body", func() {
			var err error
			request, err = http.NewRequest("POST", "/templates", strings.NewReader(`{
				"name": "a cool template",
				"text": ""
			}`))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing either template text or html"]
			}`))
		})

		It("returns a 422 when the request includes an empty html body and no text body", func() {
			var err error
			request, err = http.NewRequest("POST", "/templates", strings.NewReader(`{
				"name": "a cool template",
				"html": ""
			}`))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing either template text or html"]
			}`))
		})

		It("returns a 422 when the request includes empty text and html bodies", func() {
			var err error
			request, err = http.NewRequest("POST", "/templates", strings.NewReader(`{
				"name": "a cool template",
				"text": "",
				"html": ""
			}`))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing either template text or html"]
			}`))
		})

		It("returns a 401 when the request does not include a client id", func() {
			context.Set("client_id", "")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing client id"]
			}`))
		})

		It("returns a 409 when a template with the same name already exists", func() {
			templatesCollection.SetCall.Returns.Err = collections.DuplicateRecordError{errors.New("")}
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusConflict))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			templatesCollection.SetCall.Returns.Err = errors.New("The database is bad")
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{ "errors": ["The database is bad"] }`))
		})
	})
})
