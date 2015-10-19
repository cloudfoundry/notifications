package horde_test

import (
	"github.com/cloudfoundry-incubator/notifications/v2/horde"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("emails audience", func() {
	Describe("GenerateAudiences", func() {
		It("wraps the given list of emails in User objects", func() {
			logger := lager.NewLogger("notifications-foo")
			emails := horde.NewEmails()
			audiences, err := emails.GenerateAudiences([]string{"me@example.com"}, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(audiences).To(HaveLen(1))

			audience := audiences[0]
			Expect(audience.Users).To(Equal([]horde.User{{Email: "me@example.com"}}))
			Expect(audience.Endorsement).To(Equal("This message was sent directly to your email address."))
		})
	})
})
