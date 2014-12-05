package strategies_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailStrategy", func() {
	var emailStrategy strategies.EmailStrategy

	Describe("DispatchMail", func() {
		var fakeMailer *fakes.Mailer
		var conn *fakes.DBConn
		var options postal.Options
		var clientID string
		var emailID string
		var templatesLoader fakes.TemplatesLoader
		var scope string

		BeforeEach(func() {
			fakeMailer = fakes.NewMailer()
			templatesLoader = fakes.TemplatesLoader{}
			emailStrategy = strategies.NewEmailStrategy(fakeMailer, &templatesLoader)

			clientID = "raptors-123"
			emailID = ""

			options = postal.Options{
				Text: "email text",
				To:   "dr@strangelove.com",
			}

			conn = fakes.NewDBConn()

			templatesLoader.Templates = postal.Templates{
				Name:    "The Name",
				Subject: "the subject",
				Text:    "the text",
				HTML:    "email template",
			}
		})

		It("Calls Deliver on it's mailer with proper arguments", func() {
			emailStrategy.Dispatch(clientID, emailID, options, conn)

			users := map[string]uaa.User{"": uaa.User{Emails: []string{options.To}}}

			Expect(len(fakeMailer.DeliverArguments)).To(Equal(8))

			Expect(fakeMailer.DeliverArguments).To(ContainElement(conn))
			Expect(fakeMailer.DeliverArguments).To(ContainElement(templatesLoader.Templates))
			Expect(fakeMailer.DeliverArguments).To(ContainElement(users))
			Expect(fakeMailer.DeliverArguments).To(ContainElement(options))
			Expect(fakeMailer.DeliverArguments).To(ContainElement(cf.CloudControllerSpace{}))
			Expect(fakeMailer.DeliverArguments).To(ContainElement(cf.CloudControllerOrganization{}))
			Expect(fakeMailer.DeliverArguments).To(ContainElement(clientID))
			Expect(fakeMailer.DeliverArguments).To(ContainElement(scope))
		})
	})

	Describe("Trim", func() {
		It("Trims the recipients field", func() {
			responses, err := json.Marshal([]strategies.Response{
				{
					Status:         "delivered",
					Email:          "user@example.com",
					NotificationID: "123-456",
				},
			})

			trimmedResponses := emailStrategy.Trim(responses)

			var result []map[string]string
			err = json.Unmarshal(trimmedResponses, &result)
			if err != nil {
				panic(err)
			}

			Expect(result).To(ContainElement(map[string]string{"status": "delivered",
				"email":           "user@example.com",
				"notification_id": "123-456",
			}))
		})
	})
})
