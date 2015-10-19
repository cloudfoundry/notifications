package horde_test

import (
	"github.com/cloudfoundry-incubator/notifications/v2/horde"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("users audience", func() {
	Describe("GenerateAudiences", func() {
		It("wraps the given list of userGUIDs in User objects", func() {
			logger := lager.NewLogger("notifications-whatever")
			users := horde.NewUsers()
			audiences, err := users.GenerateAudiences([]string{"59eb64c4-728d-11e5-bf96-10ddb1aa2a2c"}, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(audiences).To(HaveLen(1))

			audience := audiences[0]
			Expect(audience.Users).To(Equal([]horde.User{{GUID: "59eb64c4-728d-11e5-bf96-10ddb1aa2a2c"}}))
			Expect(audience.Endorsement).To(Equal("This message was sent directly to you."))
		})
	})
})
