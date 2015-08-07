package templates_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/templates"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/stack"
)

var _ = Describe("GetHandler", func() {
	var (
		handler    templates.GetHandler
		context    stack.Context
		database   *fakes.Database
		writer     *httptest.ResponseRecorder
		request    *http.Request
		collection *fakes.TemplatesCollection
	)

	BeforeEach(func() {
		context = stack.NewContext()

		database = fakes.NewDatabase()
		context.Set("database", database)

		context.Set("client_id", "some-client-id")

		writer = httptest.NewRecorder()

		var err error
		request, err = http.NewRequest("GET", "/templates/some-template-id", nil)
		Expect(err).NotTo(HaveOccurred())

		collection = fakes.NewTemplatesCollection()

		handler = templates.NewGetHandler(collection)
	})

	It("gets a template", func() {
		collection.GetCall.ReturnTemplate = collections.Template{
			ID:       "some-template-id",
			Name:     "an interesting template",
			Text:     "template text",
			HTML:     "template html",
			Subject:  "template subject",
			Metadata: `{ "template": "metadata" }`,
		}

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

		Expect(collection.GetCall.TemplateID).To(Equal("some-template-id"))
		Expect(collection.GetCall.ClientID).To(Equal("some-client-id"))
		Expect(collection.GetCall.Conn).To(Equal(database.Conn))
		Expect(database.ConnectionWasCalled).To(BeTrue())
	})

	Describe("error cases", func() {
		It("responds with 404 if the collection Get returns not found", func() {
			collection.GetCall.Err = collections.NewNotFoundError("it was not found")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body).To(MatchJSON(`{"errors": ["it was not found"]}`))
		})

		It("responds with 500 if the collection Get fails", func() {
			collection.GetCall.Err = errors.New("an unknown error")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body).To(MatchJSON(`{"errors": ["an unknown error"]}`))
		})
	})
})
