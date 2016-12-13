package campaigns_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigns"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Campaign status handler", func() {
	var (
		handler                    campaigns.StatusHandler
		context                    stack.Context
		writer                     *httptest.ResponseRecorder
		request                    *http.Request
		database                   *mocks.Database
		conn                       *mocks.Connection
		campaignStatusesCollection *mocks.CampaignStatusesCollection
	)

	BeforeEach(func() {
		tokenHeader := map[string]interface{}{
			"alg": "RS256",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "some-uaa-client-id",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		token, err := jwt.Parse(helpers.BuildToken(tokenHeader, tokenClaims), func(*jwt.Token) (interface{}, error) {
			return []byte(application.UAAPublicKey), nil
		})
		Expect(err).NotTo(HaveOccurred())

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		context = stack.NewContext()
		context.Set("token", token)
		context.Set("database", database)
		context.Set("client_id", "my-client")

		writer = httptest.NewRecorder()

		campaignStatusesCollection = mocks.NewCampaignStatusesCollection()

		request, err = http.NewRequest("GET", "/campaigns/some-campaign-id/status", nil)
		Expect(err).NotTo(HaveOccurred())

		handler = campaigns.NewStatusHandler(campaignStatusesCollection)
	})

	It("gets the status of an existing campaign", func() {
		startTime, err := time.Parse(time.RFC3339, "2015-09-01T12:34:56-07:00")
		Expect(err).NotTo(HaveOccurred())

		completedTime, err := time.Parse(time.RFC3339, "2015-09-01T12:34:58-07:00")
		Expect(err).NotTo(HaveOccurred())

		campaignStatusesCollection.GetCall.Returns.CampaignStatus = collections.CampaignStatus{
			CampaignID:            "some-campaign-id",
			Status:                "completed",
			TotalMessages:         9,
			SentMessages:          6,
			QueuedMessages:        0,
			RetryMessages:         0,
			FailedMessages:        2,
			UndeliverableMessages: 1,
			StartTime:             startTime,
			CompletedTime:         &completedTime,
		}

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body).To(MatchJSON(`{
			"id": "some-campaign-id",
			"status": "completed",
			"total_messages": 9,
			"sent_messages": 6,
			"queued_messages": 0,
			"retry_messages": 0,
			"failed_messages": 2,
			"undeliverable_messages": 1,
			"start_time": "2015-09-01T12:34:56-07:00",
			"completed_time": "2015-09-01T12:34:58-07:00",
			"_links": {
				"self": {
					"href": "/campaigns/some-campaign-id/status"
				},
				"campaign": {
					"href": "/campaigns/some-campaign-id"
				}
			}
		}`))

		Expect(campaignStatusesCollection.GetCall.Receives.Connection).To(Equal(conn))
		Expect(campaignStatusesCollection.GetCall.Receives.CampaignID).To(Equal("some-campaign-id"))
		Expect(campaignStatusesCollection.GetCall.Receives.ClientID).To(Equal("my-client"))
	})

	Context("when the campaign is not yet completed", func() {
		It("returns a null completed_time value", func() {
			startTime, err := time.Parse(time.RFC3339, "2015-09-01T12:34:56-07:00")
			Expect(err).NotTo(HaveOccurred())

			campaignStatusesCollection.GetCall.Returns.CampaignStatus = collections.CampaignStatus{
				CampaignID:     "some-campaign-id",
				Status:         "sending",
				TotalMessages:  9,
				SentMessages:   5,
				QueuedMessages: 1,
				RetryMessages:  1,
				FailedMessages: 2,
				StartTime:      startTime,
				CompletedTime:  nil,
			}

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body).To(MatchJSON(`{
				"id": "some-campaign-id",
				"status": "sending",
				"total_messages": 9,
				"sent_messages": 5,
				"queued_messages": 1,
				"retry_messages": 1,
				"failed_messages": 2,
				"undeliverable_messages": 0,
				"start_time": "2015-09-01T12:34:56-07:00",
				"completed_time": null,
				"_links": {
					"self": {
						"href": "/campaigns/some-campaign-id/status"
					},
					"campaign": {
						"href": "/campaigns/some-campaign-id"
					}
				}
			}`))
		})
	})

	Context("failure cases", func() {
		It("returns a 404 with an appropriate error when the campaign statuses collection returns a not found error", func() {
			campaignStatusesCollection.GetCall.Returns.Error = collections.NotFoundError{errors.New("not found")}

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body).To(MatchJSON(`{
				"errors": [
					"not found"
				]
			}`))
		})

		It("returns a 500 with an appropriate error when the campaign statuses collection blows up", func() {
			campaignStatusesCollection.GetCall.Returns.Error = errors.New("unexpected")

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body).To(MatchJSON(`{
				"errors": [
					"unexpected"
				]
			}`))
		})
	})
})
