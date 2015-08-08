package uaa_test

import (
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenLoader", func() {
	Describe("#Load", func() {
		It("Gets a zoned client token based on hostname", func() {
			hostname := "my-uaa-zone"
			uaaFake := fakes.NewZonedUAAClient()
			uaaFake.Token = "my-fake-token"
			tokenLoader := uaa.NewTokenLoader(uaaFake)
			token, err := tokenLoader.Load(hostname)
			Expect(token).To(Equal("my-fake-token"))
			Expect(err).To(BeNil())

			Expect(uaaFake.GetClientTokenHost).To(Equal(hostname))
		})
	})
})
