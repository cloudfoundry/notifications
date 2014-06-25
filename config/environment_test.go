package config_test

import (
    "os"

    "github.com/cloudfoundry-incubator/notifications/config"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
    Describe("UAA configuration", func() {
        It("loads the values when they are set", func() {
            os.Setenv("UAA_HOST", "https://uaa.example.com")
            os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
            os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

            env := config.NewEnvironment()

            Expect(env.UAAHost).To(Equal("https://uaa.example.com"))
            Expect(env.UAAClientID).To(Equal("uaa-client-id"))
            Expect(env.UAAClientSecret).To(Equal("uaa-client-secret"))
        })

        It("panics when the values are missing", func() {
            os.Setenv("UAA_HOST", "")
            os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
            os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())

            os.Setenv("UAA_HOST", "https://uaa.example.com")
            os.Setenv("UAA_CLIENT_ID", "")
            os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())

            os.Setenv("UAA_HOST", "https://uaa.example.com")
            os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
            os.Setenv("UAA_CLIENT_SECRET", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })
})
