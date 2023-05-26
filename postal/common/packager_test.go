package common_test

import (
	"errors"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packager", func() {
	var (
		packager        common.Packager
		context         common.MessageContext
		templatesLoader *mocks.TemplatesLoader
		delivery        common.Delivery
		cloak           *mocks.Cloak
	)

	BeforeEach(func() {
		templatesLoader = mocks.NewTemplatesLoader()
		cloak = mocks.NewCloak()

		delivery = common.Delivery{
			UserGUID: "some-user-guid",
			ClientID: "some-client-id",
			Options: common.Options{
				Subject:    "Some crazy subject",
				TemplateID: "some-template-id",
				KindID:     "some-kind-id",
				HTML: common.HTML{
					BodyContent:    "<p>user supplied banana html</p>",
					BodyAttributes: "class=\"bananaBody\"",
					Head:           "<title>The title</title>",
					Doctype:        "<!DOCTYPE html>",
				},
				Text: "some-text",
			},
		}

		packager = common.NewPackager(templatesLoader, cloak)

		requestReceivedTime, _ := time.Parse(time.RFC3339Nano, "2015-06-08T14:38:03.180764129-07:00")

		context = common.MessageContext{
			From:      "banana man",
			ReplyTo:   "awesomeness",
			To:        "endless monkeys",
			Subject:   "we will be eaten",
			ClientID:  "3&3",
			MessageID: "4'4",
			Text:      "User <supplied> \"banana\" text",
			UserGUID:  "user-123",
			HTMLComponents: common.HTML{
				BodyContent:    "<p>user supplied banana html</p>",
				BodyAttributes: "class=\"bananaBody\"",
				Head:           "<title>The title</title>",
				Doctype:        "<!DOCTYPE html>",
			},
			HTML:            "<p>user supplied banana html</p>",
			Space:           "development",
			Organization:    "banana",
			TextTemplate:    "Banana preamble {{.Text}} {{.ClientID}} {{.MessageID}} {{.UserGUID}}\n{{.Endorsement}}",
			HTMLTemplate:    "<header>{{.Endorsement}}</header>\nBanana preamble {{.HTML}} {{.Text}} {{.ClientID}} {{.MessageID}} {{.UserGUID}}",
			SubjectTemplate: "The Subject: {{.Subject}}",
			Endorsement:     "This is an endorsement for the {{.Space}} space and {{.Organization}} org.",
			RequestReceived: requestReceivedTime,
		}
	})

	Describe("PrepareContext", func() {
		BeforeEach(func() {
			templatesLoader.LoadTemplatesCall.Returns.Templates = common.Templates{
				Name:    "some-name",
				Subject: "subject template: {{.Subject}}",
				Text:    "Some {{.Text}} text",
				HTML:    "<h1>{{.HTML}}</h1>",
			}
		})

		It("sets the context on the Packager", func() {
			cloak.VeilCall.Returns.CipherText = []byte("some-encrypted-text")

			var err error
			context, err = packager.PrepareContext(delivery, "some-sender@example.com", "example.com")
			Expect(err).NotTo(HaveOccurred())

			Expect(templatesLoader.LoadTemplatesCall.Receives.ClientID).To(Equal("some-client-id"))
			Expect(templatesLoader.LoadTemplatesCall.Receives.KindID).To(Equal("some-kind-id"))
			Expect(templatesLoader.LoadTemplatesCall.Receives.TemplateID).To(Equal("some-template-id"))

			Expect(cloak.VeilCall.Receives.PlainText).To(Equal([]byte("some-user-guid|some-client-id|some-kind-id")))

			Expect(context).To(Equal(common.MessageContext{
				UnsubscribeID: "some-encrypted-text",
				Domain:        "example.com",
				From:          "some-sender@example.com",
				Subject:       "Some crazy subject",
				UserGUID:      "some-user-guid",
				ClientID:      "some-client-id",
				Text:          "some-text",
				HTML:          "<p>user supplied banana html</p>",
				HTMLComponents: common.HTML{
					BodyContent:    "<p>user supplied banana html</p>",
					BodyAttributes: "class=\"bananaBody\"",
					Head:           "<title>The title</title>",
					Doctype:        "<!DOCTYPE html>",
				},
				TextTemplate:      "Some {{.Text}} text",
				HTMLTemplate:      "<h1>{{.HTML}}</h1>",
				SubjectTemplate:   "subject template: {{.Subject}}",
				KindDescription:   "some-kind-id",
				SourceDescription: "some-client-id",
			}))
		})

		Context("when the template cannot be loaded", func() {
			It("returns an error", func() {
				templatesLoader.LoadTemplatesCall.Returns.Error = errors.New("some error")

				_, err := packager.PrepareContext(delivery, "some-sender", "some-domain")
				Expect(err).To(MatchError(errors.New("some error")))
			})
		})
	})

	Describe("Pack", func() {
		It("packs a message for delivery", func() {
			msg, err := packager.Pack(context)
			Expect(err).NotTo(HaveOccurred())
			Expect(msg.From).To(Equal("banana man"))
			Expect(msg.ReplyTo).To(Equal("awesomeness"))
			Expect(msg.To).To(Equal("endless monkeys"))
			Expect(msg.Subject).To(Equal("The Subject: we will be eaten"))
			Expect(msg.Body).To(ConsistOf([]mail.Part{
				{
					ContentType: "text/plain",
					Content:     "Banana preamble User <supplied> \"banana\" text 3&3 4'4 user-123\nThis is an endorsement for the development space and banana org.",
				},
				{
					ContentType: "text/html",
					Content:     "<!DOCTYPE html>\n<head><title>The title</title></head>\n<html>\n\t<body class=\"bananaBody\">\n\t\t<header>This is an endorsement for the development space and banana org.</header>\nBanana preamble <p>user supplied banana html</p> User &lt;supplied&gt; &#34;banana&#34; text 3&amp;3 4&#39;4 user-123\n\t</body>\n</html>",
				},
			}))
			Expect(msg.Headers).To(ContainElement("X-CF-Client-ID: 3&3"))
			Expect(msg.Headers).To(ContainElement("X-CF-Notification-ID: 4'4"))
			Expect(msg.Headers).To(ContainElement("X-CF-Notification-Request-Received: 2015-06-08T14:38:03.180764129-07:00"))

			var formattedTimestamp string
			prefix := "X-CF-Notification-Timestamp: "
			for _, header := range msg.Headers {
				if strings.Contains(header, prefix) {
					formattedTimestamp = strings.TrimPrefix(header, prefix)
					break
				}
			}
			Expect(formattedTimestamp).NotTo(BeEmpty())

			timestamp, err := time.Parse(time.RFC3339Nano, formattedTimestamp)
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(BeTemporally("~", time.Now(), 2*time.Second))
		})
	})

	Describe("CompileParts", func() {
		It("returns the compiled parts containing both the plaintext and html portions, escaping variables for the html portion only", func() {
			parts, err := packager.CompileParts(context)
			if err != nil {
				panic(err)
			}

			textBody := `Banana preamble User <supplied> "banana" text 3&3 4'4 user-123
This is an endorsement for the development space and banana org.`
			htmlBody := `<!DOCTYPE html>
<head><title>The title</title></head>
<html>
	<body class="bananaBody">
		<header>This is an endorsement for the development space and banana org.</header>
Banana preamble <p>user supplied banana html</p> User &lt;supplied&gt; &#34;banana&#34; text 3&amp;3 4&#39;4 user-123
	</body>
</html>`

			Expect(parts).To(ContainElement(mail.Part{
				ContentType: "text/plain",
				Content:     textBody,
			}))
			Expect(parts).To(ContainElement(mail.Part{
				ContentType: "text/html",
				Content:     htmlBody,
			}))
		})

		Context("when no html is set", func() {
			It("only sends a plaintext of the email", func() {
				context.HTML = ""

				parts, err := packager.CompileParts(context)
				if err != nil {
					panic(err)
				}

				textBody := `Banana preamble User <supplied> "banana" text 3&3 4'4 user-123
This is an endorsement for the development space and banana org.`
				Expect(parts).To(ConsistOf([]mail.Part{
					{
						ContentType: "text/plain",
						Content:     textBody,
					},
				}))
			})
		})

		Context("when no text is set", func() {
			It("omits the plaintext portion of the email", func() {
				context.Text = ""

				parts, err := packager.CompileParts(context)
				if err != nil {
					panic(err)
				}

				htmlBody := `<!DOCTYPE html>
<head><title>The title</title></head>
<html>
	<body class="bananaBody">
		<header>This is an endorsement for the development space and banana org.</header>
Banana preamble <p>user supplied banana html</p>  3&amp;3 4&#39;4 user-123
	</body>
</html>`
				Expect(parts).To(ConsistOf([]mail.Part{
					{
						ContentType: "text/html",
						Content:     htmlBody,
					},
				}))
			})
		})
	})
})
