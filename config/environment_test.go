package config_test

import (
    "os"

    "github.com/cloudfoundry-incubator/notifications/config"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
    variables := map[string]string{
        "UAA_HOST":          os.Getenv("UAA_HOST"),
        "UAA_CLIENT_ID":     os.Getenv("UAA_CLIENT_ID"),
        "UAA_CLIENT_SECRET": os.Getenv("UAA_CLIENT_SECRET"),
        "SMTP_USER":         os.Getenv("SMTP_USER"),
        "SMTP_PASS":         os.Getenv("SMTP_PASS"),
        "SMTP_HOST":         os.Getenv("SMTP_HOST"),
        "SMTP_PORT":         os.Getenv("SMTP_PORT"),
    }

    AfterEach(func() {
        for key, value := range variables {
            os.Setenv(key, value)
        }
    })

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

    Describe("SMTP configuration", func() {
        It("loads the values when they are present", func() {
            os.Setenv("SMTP_USER", "my-smtp-user")
            os.Setenv("SMTP_PASS", "my-smtp-password")
            os.Setenv("SMTP_HOST", "smtp.example.com")
            os.Setenv("SMTP_PORT", "567")
            os.Setenv("SMTP_TLS", "true")

            env := config.NewEnvironment()

            Expect(env.SMTPUser).To(Equal("my-smtp-user"))
            Expect(env.SMTPPass).To(Equal("my-smtp-password"))
            Expect(env.SMTPHost).To(Equal("smtp.example.com"))
            Expect(env.SMTPPort).To(Equal("567"))
            Expect(env.SMTPTLS).To(BeTrue())
        })

        It("panics when SMTP_TLS is not a boolean", func() {
            os.Setenv("SMTP_TLS", "banana")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })

        It("panics when the values are missing", func() {
            os.Setenv("SMTP_USER", "my-smtp-user")
            os.Setenv("SMTP_PASS", "my-smtp-password")
            os.Setenv("SMTP_HOST", "smtp.example.com")
            os.Setenv("SMTP_PORT", "567")
            os.Setenv("SMTP_TLS", "")

            Expect(func() {
                config.NewEnvironment()
            }).NotTo(Panic())

            os.Setenv("SMTP_USER", "")
            os.Setenv("SMTP_PASS", "my-smtp-password")
            os.Setenv("SMTP_HOST", "smtp.example.com")
            os.Setenv("SMTP_PORT", "567")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())

            os.Setenv("SMTP_USER", "my-smtp-user")
            os.Setenv("SMTP_PASS", "")
            os.Setenv("SMTP_HOST", "smtp.example.com")
            os.Setenv("SMTP_PORT", "567")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())

            os.Setenv("SMTP_USER", "my-smtp-user")
            os.Setenv("SMTP_PASS", "my-smtp-password")
            os.Setenv("SMTP_HOST", "")
            os.Setenv("SMTP_PORT", "567")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())

            os.Setenv("SMTP_USER", "my-smtp-user")
            os.Setenv("SMTP_PASS", "my-smtp-password")
            os.Setenv("SMTP_HOST", "smtp.example.com")
            os.Setenv("SMTP_PORT", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })

    Describe("Sender configuration", func() {
        It("loads the SENDER environment variable when it is present", func() {
            os.Setenv("SENDER", "my-email@example.com")

            env := config.NewEnvironment()

            Expect(env.Sender).To(Equal("my-email@example.com"))
        })

        It("panics if the SENDER variable is missing", func() {
            os.Setenv("SENDER", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })
})
