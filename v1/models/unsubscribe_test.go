package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unsubscribes", func() {
	var unsubscribes models.Unsubscribes

	Describe("Contains", func() {
		BeforeEach(func() {
			unsubscribes = models.Unsubscribes([]models.Unsubscribe{
				{
					ClientID: "client-id",
					KindID:   "kind-id",
					UserID:   "user-id",
				},
				{
					ClientID: "client-id",
					KindID:   "other-kind-id",
					UserID:   "user-id",
				},
			})
		})

		Context("when the unsubscribe is in the set", func() {
			It("returns true", func() {
				Expect(unsubscribes.Contains("client-id", "kind-id")).To(BeTrue())
			})
		})

		Context("when the unsubscribe is not in the set", func() {
			It("returns false", func() {
				Expect(unsubscribes.Contains("client-id", "bad-kind-id")).To(BeFalse())
			})
		})
	})
})
