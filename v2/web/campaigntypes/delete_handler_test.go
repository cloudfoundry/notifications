package campaigntypes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
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
		campaignTypesCollection *fakes.CampaignTypesCollection
		database                *fakes.Database
	)

	BeforeEach(func() {
		context = stack.NewContext()

		database = fakes.NewDatabase()
		context.Set("database", database)

		context.Set("client_id", "some-client-id")

		writer = httptest.NewRecorder()

		var err error
		request, err = http.NewRequest("DELETE", "/senders/some-sender-id/campaign_types/some-campaign-type-id", nil)
		Expect(err).NotTo(HaveOccurred())

		campaignTypesCollection = fakes.NewCampaignTypesCollection()

		handler = campaigntypes.NewDeleteHandler(campaignTypesCollection)
	})

	It("deletes a campaign type", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(campaignTypesCollection.DeleteCall.CampaignTypeID).To(Equal("some-campaign-type-id"))
		Expect(campaignTypesCollection.DeleteCall.SenderID).To(Equal("some-sender-id"))
		Expect(campaignTypesCollection.DeleteCall.ClientID).To(Equal("some-client-id"))
		Expect(campaignTypesCollection.DeleteCall.Conn).To(Equal(database.Conn))
		Expect(writer.Body.String()).To(BeEmpty())
	})

	Context("when an error occurs", func() {
		Context("when the campaign type cannot be found", func() {
			BeforeEach(func() {
				campaignTypesCollection.DeleteCall.Err = collections.NewNotFoundError("Campaign type some-campaign-type-id not found")
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
				campaignTypesCollection.DeleteCall.Err = errors.New("nope")
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
