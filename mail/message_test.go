package mail_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/mail"

    . "github.com/onsi/ginkgo"
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
                Body:    "Banana",
            }
        })

        It("returns a populated data mail field as a string", func() {
            parts := strings.Split(msg.Data(), "\n")
            Expect(parts).To(Equal([]string{
                "From: me@example.com",
                "To: you@example.com",
                "Subject: Super Urgent! Read Now!",
                "MIME-Version: 1.0",
                "Content-Type: multipart/alternative; boundary=\"our-content-boundary\"",
                "",
                "Banana",
            }))
        })

        It("includes headers in the response if there are any", func() {
            msg.Headers = append(msg.Headers, "X-ClientID: banana")

            parts := strings.Split(msg.Data(), "\n")
            Expect(parts).To(Equal([]string{
                "X-ClientID: banana",
                "From: me@example.com",
                "To: you@example.com",
                "Subject: Super Urgent! Read Now!",
                "MIME-Version: 1.0",
                "Content-Type: multipart/alternative; boundary=\"our-content-boundary\"",
                "",
                "Banana",
            }))
        })
    })
})
