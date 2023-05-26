package v1

import (
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Getting a Message's status", func() {
	var (
		clientToken uaa.Token
		client      *support.Client
	)

	type MessageReponse struct {
		Status  int
		Message support.Message
		Error   error
	}

	BeforeEach(func() {
		clientToken = GetClientTokenFor("notification-sender")
		client = support.NewClient(Servers.Notifications.URL())
	})

	It("Gets a message's status", func() {
		var messageGUID string

		By("sending a notification to an email address", func() {
			status, responses, err := client.Notify.Email(clientToken.Access, "John User <user@example.com>", support.Notify{
				HTML:    "<header>this is an acceptance test</header>",
				Text:    "some text for the email",
				Subject: "my-special-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(responses).To(HaveLen(1))
			response := responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())

			messageGUID = response.NotificationID
		})

		By("polling the messages endpoint", func() {
			Eventually(func() MessageReponse {
				status, message, err := client.Messages.Get(clientToken.Access, messageGUID)
				return MessageReponse{
					Status:  status,
					Message: message,
					Error:   err,
				}
			}, 10*time.Second).Should(Equal(MessageReponse{
				Status:  http.StatusOK,
				Message: support.Message{Status: "delivered"},
				Error:   nil,
			}))
		})
	})
})
