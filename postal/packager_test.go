package postal_test

import (
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packager", func() {
	var (
		packager postal.Packager
		context  postal.MessageContext
		client   mail.Client
	)

	BeforeEach(func() {
		client = mail.Client{}
		html := postal.HTML{
			BodyContent:    "<p>user supplied banana html</p>",
			BodyAttributes: "class=\"bananaBody\"",
			Head:           "<title>The title</title>",
			Doctype:        "<!DOCTYPE html>",
		}

		context = postal.MessageContext{
			From:            "banana man",
			ReplyTo:         "awesomeness",
			To:              "endless monkeys",
			Subject:         "we will be eaten",
			ClientID:        "3&3",
			MessageID:       "4'4",
			Text:            "User <supplied> \"banana\" text",
			UserGUID:        "user-123",
			HTMLComponents:  html,
			HTML:            html.BodyContent,
			Space:           "development",
			Organization:    "banana",
			TextTemplate:    "Banana preamble {{.Text}} {{.ClientID}} {{.MessageID}} {{.UserGUID}}\n{{.Endorsement}}",
			HTMLTemplate:    "<header>{{.Endorsement}}</header>\nBanana preamble {{.HTML}} {{.Text}} {{.ClientID}} {{.MessageID}} {{.UserGUID}}",
			SubjectTemplate: "The Subject: {{.Subject}}",
			Endorsement:     "This is an endorsement for the {{.Space}} space and {{.Organization}} org.",
		}
		packager = postal.NewPackager()
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
				packager = postal.NewPackager()

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
				packager = postal.NewPackager()

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
