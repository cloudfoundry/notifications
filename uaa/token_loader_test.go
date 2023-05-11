package uaa_test

import (
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenLoader", func() {
	Describe("#Load", func() {
		It("Gets a zoned client token based on hostname", func() {
			uaaClient := mocks.NewZonedUAAClient()
			uaaClient.GetClientTokenCall.Returns.Token = "my-fake-token"

			tokenLoader := uaa.NewTokenLoader(uaaClient)

			token, err := tokenLoader.Load("my-uaa-zone")
			Expect(token).To(Equal("my-fake-token"))
			Expect(err).To(BeNil())

			Expect(uaaClient.GetClientTokenCall.Receives.Host).To(Equal("my-uaa-zone"))
		})
	})
})
