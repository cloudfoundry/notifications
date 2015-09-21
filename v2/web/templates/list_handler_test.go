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

var _ = Describe("ListHandler", func() {
	var (
		handler             templates.ListHandler
		context             stack.Context
		conn                *mocks.Connection
		database            *mocks.Database
		writer              *httptest.ResponseRecorder
		request             *http.Request
		templatesCollection *mocks.TemplatesCollection
	)

	BeforeEach(func() {
		context = stack.NewContext()

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn
		context.Set("database", database)

		context.Set("client_id", "some-client-id")

		writer = httptest.NewRecorder()

		var err error
		request, err = http.NewRequest("GET", "/templates/some-template-id", nil)
		Expect(err).NotTo(HaveOccurred())

		templatesCollection = mocks.NewTemplatesCollection()
		templatesCollection.ListCall.Returns.Templates = []collections.Template{
			{
				ID:       "some-template-id",
				Name:     "an interesting template",
				Text:     "template text",
				HTML:     "template html",
				Subject:  "template subject",
				Metadata: `{ "template": "metadata" }`,
			},
		}

		handler = templates.NewListHandler(templatesCollection)
	})

	It("lists all the templates", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body).To(MatchJSON(`{
			"templates": [
				{
					"id": "some-template-id",
					"name": "an interesting template",
					"text": "template text",
					"html": "template html",
					"subject": "template subject",
					"metadata": {
						"template": "metadata"
					},
					"_links": {
						"self": {
							"href": "/templates/some-template-id"
						}
					}
				}
			],
			"_links": {
				"self": {
					"href": "/templates"
				}
			}
		}`))

		Expect(templatesCollection.ListCall.Receives.Connection).To(Equal(conn))
		Expect(templatesCollection.ListCall.Receives.ClientID).To(Equal("some-client-id"))
	})

	Context("failure cases", func() {
		It("returns a 500 if an error occurs", func() {
			templatesCollection.ListCall.Returns.Error = errors.New("collection failure")

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body).To(MatchJSON(`{
				"errors": [
					"collection failure"
				]
			}`))
		})
	})
})
