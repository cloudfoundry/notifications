package mail_test

import (
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/mail"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Message", func() {
	Describe("Data", func() {
		var msg mail.Message

		BeforeEach(func() {
			msg = mail.Message{
				From:    "me@example.com",
				To:      "you@example.com",
				Subject: "Super Urgent! Read Now!",
				Body: []mail.Part{
					{
						ContentType: "text/plain",
						Content:     "Banana",
					},
					{
						ContentType: "text/html",
						Content:     "<header>banana</header>",
					},
				},
			}
		})

		It("returns a populated data mail field as a string", func() {
			parts := strings.Split(msg.Data(), "\n")
			boundary := msg.Boundary()

			Expect(parts).To(ConsistOf([]string{
				"From: me@example.com",
				"To: you@example.com",
				"Subject: Super Urgent! Read Now!",
				"Content-Type: multipart/alternative; boundary=" + boundary,
				"Date: " + time.Now().Format(time.RFC822Z),
				"Mime-Version: 1.0",
				"",
				"--" + boundary,
				"Content-Type: text/plain; charset=UTF-8",
				"Content-Transfer-Encoding: quoted-printable",
				"",
				"Banana",
				"--" + boundary,
				"Content-Type: text/html; charset=UTF-8",
				"Content-Transfer-Encoding: quoted-printable",
				"",
				"<header>banana</header>",
				"--" + boundary + "--",
				"",
			}))
		})

		Context("when optional fields are present", func() {
			It("includes Reply-To in message body", func() {
				msg.ReplyTo = "banana@chiquita.com"
				parts := strings.Split(msg.Data(), "\n")
				boundary := msg.Boundary()

				Expect(parts).To(ConsistOf([]string{
					"From: me@example.com",
					"Reply-To: banana@chiquita.com",
					"To: you@example.com",
					"Subject: Super Urgent! Read Now!",
					"Content-Type: multipart/alternative; boundary=" + boundary,
					"Date: " + time.Now().Format(time.RFC822Z),
					"Mime-Version: 1.0",
					"",
					"--" + boundary,
					"Content-Type: text/plain; charset=UTF-8",
					"Content-Transfer-Encoding: quoted-printable",
					"",
					"Banana",
					"--" + boundary,
					"Content-Type: text/html; charset=UTF-8",
					"Content-Transfer-Encoding: quoted-printable",
					"",
					"<header>banana</header>",
					"--" + boundary + "--",
					"",
				}))
			})

			It("includes headers in the response if there are any", func() {
				msg.Headers = append(msg.Headers, "X-ClientID: banana")
				parts := strings.Split(msg.Data(), "\n")
				boundary := msg.Boundary()

				Expect(parts).To(ConsistOf([]string{
					"From: me@example.com",
					"To: you@example.com",
					"Subject: Super Urgent! Read Now!",
					"X-ClientID: banana",
					"Content-Type: multipart/alternative; boundary=" + boundary,
					"Date: " + time.Now().Format(time.RFC822Z),
					"Mime-Version: 1.0",
					"",
					"--" + boundary,
					"Content-Type: text/plain; charset=UTF-8",
					"Content-Transfer-Encoding: quoted-printable",
					"",
					"Banana",
					"--" + boundary,
					"Content-Type: text/html; charset=UTF-8",
					"Content-Transfer-Encoding: quoted-printable",
					"",
					"<header>banana</header>",
					"--" + boundary + "--",
					"",
				}))
			})

			It("includes only the parts necessary", func() {
				msg.Body = []mail.Part{
					{
						ContentType: "text/html",
						Content:     "<header>banana</header>",
					},
				}

				parts := strings.Split(msg.Data(), "\n")

				Expect(parts).To(Equal([]string{
					"Date: " + time.Now().Format(time.RFC822Z),
					"Mime-Version: 1.0",
					"Content-Type: text/html; charset=UTF-8",
					"Content-Transfer-Encoding: quoted-printable",
					"From: me@example.com",
					"To: you@example.com",
					"Subject: Super Urgent! Read Now!",
					"",
					"<header>banana</header>",
				}))
			})
		})
	})
})
