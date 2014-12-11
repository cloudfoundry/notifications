package strategies_test

import (
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

			users := map[string]uaa.User{options.To: uaa.User{Emails: []string{options.To}}}

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
})
