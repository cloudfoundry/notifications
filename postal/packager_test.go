package postal_test

import (
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/postal"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Packager", func() {
    var packager postal.Packager
    var context postal.MessageContext
    var client mail.Client

    BeforeEach(func() {
        client = mail.Client{}
        context = postal.MessageContext{
            From:            "banana man",
            ReplyTo:         "awesomeness",
            To:              "endless monkeys",
            Subject:         "we will be eaten",
            ClientID:        "3&3",
            MessageID:       "4'4",
            Text:            "User <supplied> \"banana\" text",
            HTML:            "<p>user supplied banana html</p>",
            TextTemplate:    "Banana preamble {{.Text}} {{.ClientID}} {{.MessageID}}",
            HTMLTemplate:    "Banana preamble {{.HTML}} {{.Text}} {{.ClientID}} {{.MessageID}}",
            SubjectTemplate: "The Subject: {{.Subject}}",
        }
        packager = postal.NewPackager()
    })

    Describe("CompileBody", func() {
        It("returns the compiled email containing both the plaintext and html portions, escaping variables for the html portion only", func() {
            body, err := packager.CompileBody(context)
            if err != nil {
                panic(err)
            }

            emailBody := `
This is a multi-part message in MIME format...

--our-content-boundary
Content-type: text/plain

Banana preamble User <supplied> "banana" text 3&3 4'4
--our-content-boundary
Content-Type: text/html
Content-Disposition: inline
Content-Transfer-Encoding: quoted-printable

<html>
    <body>
        Banana preamble <p>user supplied banana html</p> User &lt;supplied&gt; &#34;banana&#34; text 3&amp;3 4&#39;4
    </body>
</html>
--our-content-boundary--`

            Expect(body).To(Equal(emailBody))
        })

        Context("when no html is set", func() {
            It("only sends a plaintext of the email", func() {
                context.HTML = ""
                packager = postal.NewPackager()

                body, err := packager.CompileBody(context)
                if err != nil {
                    panic(err)
                }

                emailBody := `
This is a multi-part message in MIME format...

--our-content-boundary
Content-type: text/plain

Banana preamble User <supplied> "banana" text 3&3 4'4
--our-content-boundary--`
                Expect(body).To(Equal(emailBody))
            })
        })

        Context("when no text is set", func() {
            It("omits the plaintext portion of the email", func() {
                context.Text = ""
                packager = postal.NewPackager()

                body, err := packager.CompileBody(context)
                if err != nil {
                    panic(err)
                }

                emailBody := `
This is a multi-part message in MIME format...

--our-content-boundary
Content-Type: text/html
Content-Disposition: inline
Content-Transfer-Encoding: quoted-printable

<html>
    <body>
        Banana preamble <p>user supplied banana html</p>  3&amp;3 4&#39;4
    </body>
</html>
--our-content-boundary--`
                Expect(body).To(Equal(emailBody))
            })
        })
    })

})
