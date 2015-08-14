package templates_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteHandler", func() {
	var (
		handler    templates.DeleteHandler
		writer     *httptest.ResponseRecorder
		request    *http.Request
		database   *fakes.Database
		collection *fakes.TemplatesCollection
		context    stack.Context
	)

	BeforeEach(func() {
		writer = httptest.NewRecorder()
		database = fakes.NewDatabase()

		context = stack.NewContext()
		context.Set("database", database)
		context.Set("client_id", "some-client-id")

		var err error
		request, err = http.NewRequest("DELETE", "/templates/some-template-id", nil)
		Expect(err).NotTo(HaveOccurred())

		collection = fakes.NewTemplatesCollection()
		handler = templates.NewDeleteHandler(collection)
	})

	It("deletes a template", func() {
		collection.GetCall.Returns.Template = collections.Template{
			ID: "some-template-id",
		}
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(writer.Body.String()).To(BeEmpty())

		Expect(collection.GetCall.Receives.Connection).To(Equal(database.Connection()))
		Expect(collection.GetCall.Receives.TemplateID).To(Equal("some-template-id"))
		Expect(collection.GetCall.Receives.ClientID).To(Equal("some-client-id"))

		Expect(collection.DeleteCall.Receives.Connection).To(Equal(database.Connection()))
		Expect(collection.DeleteCall.Receives.TemplateID).To(Equal("some-template-id"))
	})

	Context("failure cases", func() {
		It("returns a 404 if the template does not exist when getting", func() {
			collection.GetCall.Returns.Err = collections.NotFoundError{errors.New(`Template with id "some-template-id" could not be found`)}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["Template with id \"some-template-id\" could not be found"]
			}`))
		})

		It("returns a 404 if the template does not exist when deleting", func() {
			collection.DeleteCall.Returns.Err = collections.NotFoundError{errors.New(`Template with id "some-template-id" could not be found`)}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["Template with id \"some-template-id\" could not be found"]
			}`))
		})

		It("returns a 500 when the get collection call results in an unknown error", func() {
			collection.GetCall.Returns.Err = errors.New("something bad happened")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["something bad happened"]
			}`))
		})

		It("returns a 500 when the delete collection call results in an unknown error", func() {
			collection.DeleteCall.Returns.Err = errors.New("something bad happened")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["something bad happened"]
			}`))
		})
	})
})
