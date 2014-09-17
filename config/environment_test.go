package config_test

import (
    "os"

    "github.com/cloudfoundry-incubator/notifications/config"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
    var variables = map[string]string{}
    var envVars = []string{
        "CC_HOST",
        "CORS_ORIGIN",
        "DATABASE_URL",
        "DB_LOGGING_ENABLED",
        "PORT",
        "ROOT_PATH",
        "SENDER",
        "SMTP_HOST",
        "SMTP_PASS",
        "SMTP_PORT",
        "SMTP_USER",
        "TEST_MODE",
        "UAA_CLIENT_ID",
        "UAA_CLIENT_SECRET",
        "UAA_HOST",
        "VCAP_APPLICATION",
        "VERIFY_SSL",
    }

    BeforeEach(func() {
        for _, envVar := range envVars {
            variables[envVar] = os.Getenv(envVar)
        }
    })

    AfterEach(func() {
        for key, value := range variables {
            os.Setenv(key, value)
        }
    })

    Describe("Database URL", func() {
        Context("when DATABASE_URL is properly formatted", func() {
            It("converts the DATABASE_URL into a database driver DSN format", func() {
                os.Setenv("DATABASE_URL", "user-123:mypassword@example.com/banana")
                env := config.NewEnvironment()
                Expect(env.DatabaseURL).To(Equal("user-123:mypassword@tcp(example.com)/banana?parseTime=true"))
            })

            It("converts the DATABASE_URL into a database driver DSN format", func() {
                os.Setenv("DATABASE_URL", "https://user-123:mypassword@example.com/banana")
                env := config.NewEnvironment()
                Expect(env.DatabaseURL).To(Equal("user-123:mypassword@tcp(example.com)/banana?parseTime=true"))
            })
        })

        Context("when DATABASE_URL is not properly formatted", func() {
            It("panics when the url is not set", func() {
                os.Setenv("DATABASE_URL", "")
                Expect(func() {
                    config.NewEnvironment()
                }).To(Panic())
            })

            It("panics when the url is not properly formatted", func() {
                os.Setenv("DATABASE_URL", "s%%oe\\mthing!!")
                Expect(func() {
                    config.NewEnvironment()
                }).To(Panic())
            })
        })
    })

    Describe("Port configuration", func() {
        It("loads the value when it is set", func() {
            os.Setenv("PORT", "5001")
            env := config.NewEnvironment()
            Expect(env.Port).To(Equal("5001"))
        })

        It("sets the value to 3000 when it is not set", func() {
            os.Setenv("PORT", "")
            env := config.NewEnvironment()
            Expect(env.Port).To(Equal("3000"))
        })
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

        It("defaults to true when SMTP_TLS is not a boolean", func() {
            os.Setenv("SMTP_TLS", "")

            env := config.NewEnvironment()

            Expect(env.SMTPTLS).To(BeTrue())
        })

        It("does not panic when SMTP_USER and/or SMTP_PASS are empty", func() {
            os.Setenv("SMTP_USER", "")
            os.Setenv("SMTP_PASS", "")

            Expect(func() {
                config.NewEnvironment()
            }).NotTo(Panic())
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

    Describe("CloudController configuration", func() {
        It("loads the values when they are present", func() {
            os.Setenv("CC_HOST", "https://api.example.com")

            env := config.NewEnvironment()

            Expect(env.CCHost).To(Equal("https://api.example.com"))
        })

        It("panics when any of the values are missing", func() {
            os.Setenv("CC_HOST", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })

    Describe("SSL verification configuration", func() {
        It("set the value to true by default", func() {
            os.Setenv("VERIFY_SSL", "")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeTrue())
        })

        It("can be set to false", func() {
            os.Setenv("VERIFY_SSL", "false")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeFalse())
        })

        It("can be set to true", func() {
            os.Setenv("VERIFY_SSL", "TRUE")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeTrue())
        })

        It("sets the value to true if the value is non-boolean", func() {
            os.Setenv("VERIFY_SSL", "")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeTrue())
        })
    })

    Describe("RootPath config", func() {
        It("loads the config value", func() {
            os.Setenv("ROOT_PATH", "bananaDAMAGE")
            env := config.NewEnvironment()

            Expect(env.RootPath).To(Equal("bananaDAMAGE"))
        })
    })

    Describe("TestMode config", func() {
        It("sets the value to false by default", func() {
            os.Setenv("TEST_MODE", "")

            env := config.NewEnvironment()

            Expect(env.TestMode).To(BeFalse())
        })

        It("can be set to true", func() {
            os.Setenv("TEST_MODE", "TRUE")

            env := config.NewEnvironment()

            Expect(env.TestMode).To(BeTrue())
        })
    })

    Describe("InstanceIndex config", func() {
        It("sets the value if it is available", func() {
            os.Setenv("VCAP_APPLICATION", `{"instance_index":1}`)

            env := config.NewEnvironment()

            Expect(env.VCAPApplication.InstanceIndex).To(Equal(1))
        })

        It("panics if it cannot find the value", func() {
            os.Setenv("VCAP_APPLICATION", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })

    Describe("Database logging config", func() {
        It("defaults to false", func() {
            os.Setenv("DB_LOGGING_ENABLED", "")
            env := config.NewEnvironment()
            Expect(env.DBLoggingEnabled).To(BeFalse())
        })

        It("can be set to true", func() {
            os.Setenv("DB_LOGGING_ENABLED", "true")
            env := config.NewEnvironment()
            Expect(env.DBLoggingEnabled).To(BeTrue())
        })
    })

    Describe("CORS origin", func() {
        It("defaults to *", func() {
            env := config.NewEnvironment()
            Expect(env.CORSOrigin).To(Equal("*"))
        })

        It("uses the value set by CORS_ORIGIN", func() {
            os.Setenv("CORS_ORIGIN", "banana")
            env := config.NewEnvironment()
            Expect(env.CORSOrigin).To(Equal("banana"))
        })
    })
})
