package network_test

import (
	"github.com/pivotal-cf-experimental/warrant/internal/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestAuthorization", func() {
	Describe("TokenAuthorization", func() {
		It("returns a bearer token given a token value", func() {
			auth := network.NewTokenAuthorization("TOKEN")

			Expect(auth.Authorization()).To(Equal("Bearer TOKEN"))
		})
	})

	Describe("BasicAuthorization", func() {
		It("returns a basic auth header given a username and password", func() {
			auth := network.NewBasicAuthorization("username", "password")

			Expect(auth.Authorization()).To(Equal("Basic dXNlcm5hbWU6cGFzc3dvcmQ="))
		})
	})
})
