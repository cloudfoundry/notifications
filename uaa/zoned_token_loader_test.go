package uaa_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ZonedTokenLoader", func() {
	Describe("#Load", func() {
		It("Gets a zoned client token based on hostname", func() {
			hostname := "my-uaa-zone"
			uaaFake := fakes.NewZonedUAAClient()
			uaaFake.ZonedToken = "my-fake-token"
			zonedTokenLoader := uaa.NewZonedTokenLoader(uaaFake)
			token, err := zonedTokenLoader.Load(hostname)
			Expect(token).To(Equal("my-fake-token"))
			Expect(err).To(BeNil())

			Expect(uaaFake.ZonedGetClientTokenHost).To(Equal(hostname))
		})
	})
})
