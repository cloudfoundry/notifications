package handlers_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUserParams", func() {
    Describe("NewNotifyUserParams", func() {
        It("parses the body of the given request", func() {
            body := strings.NewReader(`{
                "kind": "test_email",
                "kind_description": "Descriptive Email Name",
                "source_description": "Descriptive Component Name",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`)

            params := handlers.NewNotifyUserParams(body)

            Expect(params.Kind).To(Equal("test_email"))
            Expect(params.KindDescription).To(Equal("Descriptive Email Name"))
            Expect(params.SourceDescription).To(Equal("Descriptive Component Name"))
            Expect(params.Subject).To(Equal("Summary of contents"))
            Expect(params.Text).To(Equal("Contents of the email message"))
        })

        It("does not blow up if the request body is empty", func() {
            body := strings.NewReader("")

            Expect(func() {
                handlers.NewNotifyUserParams(body)
            }).NotTo(Panic())
        })
    })

    Describe("Validate", func() {
        It("validates the required parameters in the request body", func() {
            body := strings.NewReader(`{
                "kind": "test_email",
                "kind_description": "Descriptive Email Name",
                "source_description": "Descriptive Component Name",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`)
            params := handlers.NewNotifyUserParams(body)

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))

            params.Kind = ""

            Expect(params.Validate()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(1))
            Expect(params.Errors).To(ContainElement(`"kind" is a required field`))

            params.Text = ""

            Expect(params.Validate()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(2))
            Expect(params.Errors).To(ContainElement(`"kind" is a required field`))
            Expect(params.Errors).To(ContainElement(`"text" is a required field`))

            params.Kind = "something"
            params.Text = "banana"

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))
        })
    })
})
