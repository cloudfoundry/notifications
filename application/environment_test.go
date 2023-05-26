package application_test

import (
	"errors"
	"os"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/ryanmoran/viron"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	var variables = map[string]string{}
	var envVars = []string{
		"CC_HOST",
		"CORS_ORIGIN",
		"DATABASE_URL",
		"DB_LOGGING_ENABLED",
		"DB_MAX_OPEN_CONNS",
		"DEFAULT_UAA_SCOPES",
		"DOMAIN",
		"ENCRYPTION_KEY",
		"GOBBLE_WAIT_MAX_DURATION",
		"PORT",
		"ROOT_PATH",
		"SENDER",
		"SMTP_AUTH_MECHANISM",
		"SMTP_CRAMMD5_SECRET",
		"SMTP_HOST",
		"SMTP_LOGGING_ENABLED",
		"SMTP_PASS",
		"SMTP_PORT",
		"SMTP_USER",
		"TEST_MODE",
		"UAA_CLIENT_ID",
		"UAA_CLIENT_SECRET",
		"UAA_HOST",
		"VCAP_APPLICATION",
		"VERIFY_SSL",
		"DATABASE_ENABLE_IDENTITY_VERIFICATION",
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

	Context("when an environment error occurs", func() {
		It("adds a helpful message about using the bosh release to the error message", func() {
			err := application.EnvironmentError{Err: errors.New("something is misconfigured")}
			Expect(err.Error()).To(Equal("something is misconfigured (Please see https://github.com/cloudfoundry/notifications-release to find a packaged version of notifications and see the required configuration)"))
		})
	})

	Describe("Database URL", func() {
		Context("when DATABASE_URL is properly formatted", func() {
			It("converts the DATABASE_URL into a database driver DSN format", func() {
				os.Setenv("DATABASE_URL", "user-123:mypassword@example.com/banana")
				env, err := application.NewEnvironment()
				Expect(err).NotTo(HaveOccurred())
				Expect(env.DatabaseURL).To(Equal("user-123:mypassword@tcp(example.com)/banana?parseTime=true"))
			})

			It("converts the DATABASE_URL into a database driver DSN format", func() {
				os.Setenv("DATABASE_URL", "https://user-123:mypassword@example.com/banana")
				env, err := application.NewEnvironment()
				Expect(err).NotTo(HaveOccurred())
				Expect(env.DatabaseURL).To(Equal("user-123:mypassword@tcp(example.com)/banana?parseTime=true"))
			})
		})

		Context("when DATABASE_URL is not properly formatted", func() {
			It("errors when the url is not set", func() {
				os.Setenv("DATABASE_URL", "")

				_, err := application.NewEnvironment()
				Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "DATABASE_URL"}}))
			})

			It("errors when the url is not properly formatted", func() {
				os.Setenv("DATABASE_URL", "s%%oe\\mthing!!")

				_, err := application.NewEnvironment()
				Expect(err).To(MatchError(application.EnvironmentError{Err: errors.New("Could not parse DATABASE_URL \"s%%oe\\\\mthing!!\", it does not fit format \"tcp://user:pass@host/dname\"")}))
			})
		})
	})

	Describe("Database max open connections", func() {
		It("defaults to 0", func() {
			os.Setenv("DB_MAX_OPEN_CONNS", "")
			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.DBMaxOpenConns).To(Equal(0))
		})

		It("can be configured", func() {
			os.Setenv("DB_MAX_OPEN_CONNS", "15")
			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.DBMaxOpenConns).To(Equal(15))
		})
	})

	Describe("Notifications Migrations Path", func() {
		It("infers the right location", func() {
			os.Setenv("ROOT_PATH", "/tmp/foo")
			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.ModelMigrationsPath).To(Equal("/tmp/foo/db/migrations"))
		})
	})

	Describe("Gobble Migrations Path", func() {
		It("infers the right location", func() {
			os.Setenv("ROOT_PATH", "/tmp/foo")
			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.GobbleMigrationsPath).To(Equal("/tmp/foo/gobble/migrations"))
		})
	})

	Describe("Port configuration", func() {
		It("loads the value when it is set", func() {
			os.Setenv("PORT", "5001")
			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Port).To(Equal(5001))
		})

		It("sets the value to 3000 when it is not set", func() {
			os.Setenv("PORT", "")
			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Port).To(Equal(3000))
		})
	})

	Describe("UAA configuration", func() {
		It("loads the values when they are set", func() {
			os.Setenv("UAA_HOST", "https://uaa.example.com")
			os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
			os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			Expect(env.UAAHost).To(Equal("https://uaa.example.com"))
			Expect(env.UAAClientID).To(Equal("uaa-client-id"))
			Expect(env.UAAClientSecret).To(Equal("uaa-client-secret"))
		})

		It("errors when the values are missing", func() {
			os.Setenv("UAA_HOST", "")
			os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
			os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

			_, err := application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "UAA_HOST"}}))

			os.Setenv("UAA_HOST", "https://uaa.example.com")
			os.Setenv("UAA_CLIENT_ID", "")
			os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

			_, err = application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "UAA_CLIENT_ID"}}))

			os.Setenv("UAA_HOST", "https://uaa.example.com")
			os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
			os.Setenv("UAA_CLIENT_SECRET", "")

			_, err = application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "UAA_CLIENT_SECRET"}}))
		})
	})

	Describe("SMTP configuration", func() {
		It("loads the values when they are present", func() {
			os.Setenv("SMTP_USER", "my-smtp-user")
			os.Setenv("SMTP_PASS", "my-smtp-password")
			os.Setenv("SMTP_HOST", "smtp.example.com")
			os.Setenv("SMTP_PORT", "567")
			os.Setenv("SMTP_TLS", "true")
			os.Setenv("SMTP_CRAMMD5_SECRET", "supersecret")
			os.Setenv("SMTP_AUTH_MECHANISM", "plain")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			Expect(env.SMTPUser).To(Equal("my-smtp-user"))
			Expect(env.SMTPPass).To(Equal("my-smtp-password"))
			Expect(env.SMTPHost).To(Equal("smtp.example.com"))
			Expect(env.SMTPPort).To(Equal("567"))
			Expect(env.SMTPCRAMMD5Secret).To(Equal("supersecret"))
			Expect(env.SMTPAuthMechanism).To(Equal("plain"))
			Expect(env.SMTPTLS).To(BeTrue())
		})

		It("defaults to true when SMTP_TLS is not a boolean", func() {
			os.Setenv("SMTP_TLS", "")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.SMTPTLS).To(BeTrue())
		})

		It("does not error when SMTP_USER and/or SMTP_PASS are empty", func() {
			os.Setenv("SMTP_USER", "")
			os.Setenv("SMTP_PASS", "")

			_, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
		})

		It("it errors if SMTP_AUTH_MECHANISM is not one of the three supported types", func() {
			os.Setenv("SMTP_AUTH_MECHANISM", "cram-md5")
			_, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			os.Setenv("SMTP_AUTH_MECHANISM", "plain")
			_, err = application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			os.Setenv("SMTP_AUTH_MECHANISM", "none")
			_, err = application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			os.Setenv("SMTP_AUTH_MECHANISM", "banana")
			_, err = application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: errors.New("Could not parse SMTP_AUTH_MECHANISM \"banana\", it is not one of the allowed values: [none plain cram-md5]")}))
		})

		It("errors when the values are missing", func() {
			os.Setenv("SMTP_HOST", "smtp.example.com")
			os.Setenv("SMTP_PORT", "567")
			os.Setenv("SMTP_AUTH_MECHANISM", "plain")

			_, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			os.Setenv("SMTP_HOST", "")

			_, err = application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "SMTP_HOST"}}))

			os.Setenv("SMTP_HOST", "smtp.example.com")
			os.Setenv("SMTP_PORT", "")

			_, err = application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "SMTP_PORT"}}))

			os.Setenv("SMTP_AUTH_MECHANISM", "")
			os.Setenv("SMTP_PORT", "567")

			_, err = application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "SMTP_AUTH_MECHANISM"}}))
		})
	})

	Describe("SMTP logging", func() {
		It("loads the SMTP_LOGGING_ENABLED variable when it is present", func() {
			os.Setenv("SMTP_LOGGING_ENABLED", "true")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.SMTPLoggingEnabled).To(BeTrue())
		})

		It("defaults the SMTP_LOGGING_ENABLED variable to false when it is not set", func() {
			os.Setenv("SMTP_LOGGING_ENABLED", "")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.SMTPLoggingEnabled).To(BeFalse())
		})
	})

	Describe("Sender configuration", func() {
		It("loads the SENDER environment variable when it is present", func() {
			os.Setenv("SENDER", "my-email@example.com")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Sender).To(Equal("my-email@example.com"))
		})

		It("errors if the SENDER variable is missing", func() {
			os.Setenv("SENDER", "")

			_, err := application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "SENDER"}}))
		})
	})

	Describe("CloudController configuration", func() {
		It("loads the values when they are present", func() {
			os.Setenv("CC_HOST", "https://api.example.com")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.CCHost).To(Equal("https://api.example.com"))
		})

		It("errors when any of the values are missing", func() {
			os.Setenv("CC_HOST", "")

			_, err := application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "CC_HOST"}}))
		})
	})

	Describe("SSL verification configuration", func() {
		It("set the value to true by default", func() {
			os.Setenv("VERIFY_SSL", "")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.VerifySSL).To(BeTrue())
		})

		It("can be set to false", func() {
			os.Setenv("VERIFY_SSL", "false")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.VerifySSL).To(BeFalse())
		})

		It("can be set to true", func() {
			os.Setenv("VERIFY_SSL", "TRUE")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.VerifySSL).To(BeTrue())
		})

		It("sets the value to true if the value is non-boolean", func() {
			os.Setenv("VERIFY_SSL", "")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.VerifySSL).To(BeTrue())
		})
	})

	Describe("RootPath config", func() {
		It("loads the config value", func() {
			os.Setenv("ROOT_PATH", "bananaDAMAGE")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.RootPath).To(Equal("bananaDAMAGE"))
		})

		It("expands the path when needed", func() {
			os.Setenv("HOME", "bananaDAMAGE")
			os.Setenv("ROOT_PATH", "$HOME")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.RootPath).To(Equal("bananaDAMAGE"))
		})
	})

	Describe("TestMode config", func() {
		It("sets the value to false by default", func() {
			os.Setenv("TEST_MODE", "")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.TestMode).To(BeFalse())
		})

		It("can be set to true", func() {
			os.Setenv("TEST_MODE", "TRUE")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.TestMode).To(BeTrue())
		})
	})

	Describe("InstanceIndex config", func() {
		It("sets the value if it is available", func() {
			os.Setenv("VCAP_APPLICATION", `{"instance_index":1}`)

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.VCAPApplication.InstanceIndex).To(Equal(1))
		})

		It("errors if it cannot find the value", func() {
			os.Setenv("VCAP_APPLICATION", "")

			_, err := application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "VCAP_APPLICATION"}}))
		})
	})

	Describe("Database logging config", func() {
		It("defaults to false", func() {
			os.Setenv("DB_LOGGING_ENABLED", "")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.DBLoggingEnabled).To(BeFalse())
		})

		It("can be set to true", func() {
			os.Setenv("DB_LOGGING_ENABLED", "true")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.DBLoggingEnabled).To(BeTrue())
		})
	})

	Describe("CORS origin", func() {
		It("defaults to *", func() {
			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.CORSOrigin).To(Equal("*"))
		})

		It("uses the value set by CORS_ORIGIN", func() {
			os.Setenv("CORS_ORIGIN", "banana")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.CORSOrigin).To(Equal("banana"))
		})
	})

	Describe("EncryptionKey", func() {
		It("sets the EncryptionKey if it is valid", func() {
			key := "this is a very very secret secret!!"
			os.Setenv("ENCRYPTION_KEY", key)

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.EncryptionKey).To(Equal([]byte(key)))
		})

		It("errors if it is not set", func() {
			os.Setenv("ENCRYPTION_KEY", "")

			_, err := application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "ENCRYPTION_KEY"}}))
		})
	})

	Describe("Gobble WaitMaxDuration", func() {
		It("sets the value if present", func() {
			os.Setenv("GOBBLE_WAIT_MAX_DURATION", "2500")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.GobbleWaitMaxDuration).To(Equal(2500))
		})

		It("defaults to 5000", func() {
			os.Setenv("GOBBLE_WAIT_MAX_DURATION", "")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.GobbleWaitMaxDuration).To(Equal(5000))
		})
	})

	Describe("Default UAA scopes", func() {
		It("sets the value if present", func() {
			os.Setenv("DEFAULT_UAA_SCOPES", "my-scope,banana,foo,bar")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.DefaultUAAScopes).To(Equal([]string{
				"my-scope",
				"banana",
				"foo",
				"bar",
			}))
		})
	})

	Describe("Domain", func() {
		It("sets the Domain", func() {
			os.Setenv("DOMAIN", "example.com")

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Domain).To(Equal("example.com"))
		})

		It("errors if it is not set", func() {
			os.Setenv("DOMAIN", "")

			_, err := application.NewEnvironment()
			Expect(err).To(MatchError(application.EnvironmentError{Err: viron.RequiredFieldError{Name: "DOMAIN"}}))
		})
	})
})
