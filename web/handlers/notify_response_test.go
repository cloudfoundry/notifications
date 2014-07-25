package handlers_test

import (
    "bytes"
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyResponseGenerator", func() {

    var notifyResponse handlers.NotifyResponseGenerator

    Describe("SendMailToUser", func() {

        var logger *log.Logger
        var buffer *bytes.Buffer
        var mailClient FakeMailClient

        BeforeEach(func() {
            buffer = bytes.NewBuffer([]byte{})
            logger = log.New(buffer, "", 0)
        })

        It("logs the email address of the recipient and returns the status", func() {
            messageContext := handlers.MessageContext{
                To: "fake-user@example.com",
            }

            mailClient = FakeMailClient{}

            status := notifyResponse.SendMailToUser(messageContext, logger, &mailClient)

            Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
            Expect(status).To(Equal("delivered"))
        })

        It("logs the message envelope", func() {
            messageContext := handlers.MessageContext{
                To:                     "fake-user@example.com",
                From:                   "from@email.com",
                Subject:                "the subject",
                Text:                   "body content",
                KindDescription:        "the kind description",
                PlainTextEmailTemplate: "{{.Text}}",
                SubjectEmailTemplate:   "{{.Subject}}",
            }

            mailClient = FakeMailClient{}

            notifyResponse.SendMailToUser(messageContext, logger, &mailClient)

            data := []string{
                "From: from@email.com",
                "To: fake-user@example.com",
                "Subject: the subject",
                `body content`,
            }
            results := strings.Split(buffer.String(), "\n")
            for _, item := range data {
                Expect(results).To(ContainElement(item))
            }
        })
    })

    Describe("LoadSubjectTemplate", func() {
        var manager handlers.EmailTemplateManager

        Context("when subject is not set in the params", func() {
            It("returns the subject.missing template", func() {
                manager.ReadFile = func(path string) (string, error) {
                    if strings.Contains(path, "missing") {
                        return "the missing subject", nil
                    }
                    return "incorrect", nil
                }

                manager.FileExists = func(path string) bool {
                    return false
                }

                subject := ""

                subjectTemplate, err := notifyResponse.LoadSubjectTemplate(subject, manager)
                if err != nil {
                    panic(err)
                }

                Expect(subjectTemplate).To(Equal("the missing subject"))
            })
        })

        Context("when subject is set in the params", func() {
            It("returns the subject.provided template", func() {
                manager.ReadFile = func(path string) (string, error) {
                    if strings.Contains(path, "provided") {
                        return "the provided subject", nil
                    }
                    return "incorrect", nil
                }

                manager.FileExists = func(path string) bool {
                    return false
                }

                subject := "is provided"

                subjectTemplate, err := notifyResponse.LoadSubjectTemplate(subject, manager)
                if err != nil {
                    panic(err)
                }

                Expect(subjectTemplate).To(Equal("the provided subject"))
            })
        })
    })

    Describe("LoadBodyTemplates", func() {
        var manager handlers.EmailTemplateManager

        Context("loadSpace is true", func() {
            It("returns the space templates", func() {

                manager.ReadFile = func(path string) (string, error) {
                    if strings.Contains(path, "space") && strings.Contains(path, "text") {
                        return "space plain text", nil
                    }
                    if strings.Contains(path, "space") && strings.Contains(path, "html") {
                        return "space html code", nil
                    }
                    return "incorrect", nil
                }

                manager.FileExists = func(path string) bool {
                    return false
                }

                notifyResponse := handlers.NotifyResponseGenerator{}

                plain, html, err := notifyResponse.LoadBodyTemplates(true, manager)
                if err != nil {
                    panic(err)
                }

                Expect(plain).To(Equal("space plain text"))
                Expect(html).To(Equal("space html code"))
            })
        })

        Context("loadSpace is false", func() {
            It("returns the user templates", func() {
                manager.ReadFile = func(path string) (string, error) {
                    if strings.Contains(path, "user") && strings.Contains(path, "text") {
                        return "user plain text", nil
                    }
                    if strings.Contains(path, "user") && strings.Contains(path, "html") {
                        return "user html code", nil
                    }
                    return "incorrect", nil
                }

                manager.FileExists = func(path string) bool {
                    return false
                }

                notifyResponse := handlers.NotifyResponseGenerator{}

                plain, html, err := notifyResponse.LoadBodyTemplates(false, manager)
                if err != nil {
                    panic(err)
                }

                Expect(plain).To(Equal("user plain text"))
                Expect(html).To(Equal("user html code"))
            })
        })
    })
})
