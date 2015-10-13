package campaigntypes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigntypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/stack"
)

var _ = Describe("DeleteHandler", func() {
	var (
		context                 stack.Context
		writer                  *httptest.ResponseRecorder
		request                 *http.Request
		handler                 campaigntypes.DeleteHandler
		campaignTypesCollection *mocks.CampaignTypesCollection
		database                *mocks.Database
		conn                    *mocks.Connection
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
		request, err = http.NewRequest("DELETE", "/campaign_types/some-campaign-type-id", nil)
		Expect(err).NotTo(HaveOccurred())

		campaignTypesCollection = mocks.NewCampaignTypesCollection()

		handler = campaigntypes.NewDeleteHandler(campaignTypesCollection)
	})

	It("deletes a campaign type", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(campaignTypesCollection.DeleteCall.Receives.CampaignTypeID).To(Equal("some-campaign-type-id"))
		Expect(campaignTypesCollection.DeleteCall.Receives.ClientID).To(Equal("some-client-id"))
		Expect(campaignTypesCollection.DeleteCall.Receives.Conn).To(Equal(conn))
		Expect(writer.Body.String()).To(BeEmpty())
	})

	Context("when an error occurs", func() {
		Context("when the campaign type cannot be found", func() {
			BeforeEach(func() {
				campaignTypesCollection.DeleteCall.Returns.Err = collections.NotFoundError{errors.New("Campaign type some-campaign-type-id not found")}
			})

			It("returns a 404 and the error", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusNotFound))
				Expect(writer.Body.String()).To(MatchJSON(`{
					"errors": [
						"Campaign type some-campaign-type-id not found"
					]
				}`))
			})
		})

		Context("when the collection delete call returns an error", func() {
			BeforeEach(func() {
				campaignTypesCollection.DeleteCall.Returns.Err = errors.New("nope")
			})

			It("returns a 500 and the error", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusInternalServerError))
				Expect(writer.Body.String()).To(MatchJSON(`{
					"errors": [
						"Delete failed with error: nope"
					]
				}`))
			})
		})
	})
})
