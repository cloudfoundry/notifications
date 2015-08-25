package templates_test

import (
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

var _ = Describe("GetHandler", func() {
	var (
		handler    templates.GetHandler
		context    stack.Context
		database   *mocks.Database
		writer     *httptest.ResponseRecorder
		request    *http.Request
		collection *mocks.TemplatesCollection
	)

	BeforeEach(func() {
		context = stack.NewContext()

		database = mocks.NewDatabase()
		context.Set("database", database)

		context.Set("client_id", "some-client-id")

		writer = httptest.NewRecorder()

		var err error
		request, err = http.NewRequest("GET", "/templates/some-template-id", nil)
		Expect(err).NotTo(HaveOccurred())

		collection = mocks.NewTemplatesCollection()

		handler = templates.NewGetHandler(collection)
	})

	It("gets a template", func() {
		collection.GetCall.Returns.Template = collections.Template{
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

		Expect(collection.GetCall.Receives.TemplateID).To(Equal("some-template-id"))
		Expect(collection.GetCall.Receives.ClientID).To(Equal("some-client-id"))
		Expect(collection.GetCall.Receives.Connection).To(Equal(database.Conn))
		Expect(database.ConnectionWasCalled).To(BeTrue())
	})

	Describe("error cases", func() {
		It("responds with 404 if the collection Get returns not found", func() {
			collection.GetCall.Returns.Error = collections.NotFoundError{errors.New("it was not found")}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body).To(MatchJSON(`{"errors": ["it was not found"]}`))
		})

		It("responds with 500 if the collection Get fails", func() {
			collection.GetCall.Returns.Error = errors.New("an unknown error")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body).To(MatchJSON(`{"errors": ["an unknown error"]}`))
		})
	})
})
