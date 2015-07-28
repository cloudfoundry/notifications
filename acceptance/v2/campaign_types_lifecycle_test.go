package v2

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/v2/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Campaign types lifecycle", func() {
	var (
		client *support.Client
		token  uaa.Token
		sender support.Sender
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host: Servers.Notifications.URL(),
		})
		token = GetClientTokenFor("my-client", "uaa")
		var err error
		sender, err = client.Senders.Create("my-sender", token.Access)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can create and show a new campaign type", func() {
		var campaignType support.CampaignType
		var err error
		By("creating a campaign type", func() {
			campaignType, err = client.CampaignTypes.Create(sender.ID, "some-campaign-type", "a great campaign type", "", false, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(campaignType.Name).To(Equal("some-campaign-type"))
			Expect(campaignType.Description).To(Equal("a great campaign type"))
			Expect(campaignType.Critical).To(BeFalse())
			Expect(campaignType.TemplateID).To(BeEmpty())
		})

		By("creating it again with the same name", func() {
			campaignType, err = client.CampaignTypes.Create(sender.ID, "some-campaign-type", "another great campaign type", "", false, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(campaignType.Name).To(Equal("some-campaign-type"))
			Expect(campaignType.Description).To(Equal("a great campaign type"))
		})

		By("showing the newly created campaign type", func() {
			gottenCampaignType, err := client.CampaignTypes.Show(sender.ID, campaignType.ID, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(gottenCampaignType.Name).To(Equal("some-campaign-type"))
			Expect(gottenCampaignType.Description).To(Equal("a great campaign type"))
		})
	})

	It("does not leak the existence of unauthorized campaign types", func() {
		var campaignType support.CampaignType
		var otherSender support.Sender
		var err error

		By("creating a campaign type belonging to 'my-sender'", func() {
			campaignType, err = client.CampaignTypes.Create(sender.ID, "some-campaign-type", "a great campaign type", "", false, token.Access)
			Expect(err).NotTo(HaveOccurred())
		})

		By("creating a sender that is not 'my-sender'", func() {
			otherSender, err = client.Senders.Create("some-other-sender", token.Access)
			Expect(err).NotTo(HaveOccurred())
		})

		By("verifying that you cannot get a campaign type belonging to a different sender", func() {
			_, err := client.CampaignTypes.Show(otherSender.ID, campaignType.ID, token.Access)

			Expect(err.(support.NotFoundError).Status).To(Equal(http.StatusNotFound))
			expectedErrorMessage := fmt.Sprintf("{\"error\": \"campaign type %s not found\"}", campaignType.ID)
			Expect(err.(support.NotFoundError).Body).To(MatchJSON(expectedErrorMessage))
		})
	})
})
