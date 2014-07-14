package handlers_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifySpaceParams", func() {
    Describe("NewNotifySpaceParams", func() {
        It("parses the request body into the struct", func() {
            body := strings.NewReader(`{"kind": "the kind", "text": "the text"}`)
            params := handlers.NewNotifySpaceParams(body)

            Expect(params.Kind).To(Equal("the kind"))
            Expect(params.Text).To(Equal("the text"))
        })
    })

    Describe("Validate", func() {
        It("returns true if request is valid", func() {
            body := strings.NewReader(`{"kind": "the kind", "text": "the text"}`)
            params := handlers.NewNotifySpaceParams(body)

            Expect(params.Validate()).To(BeTrue())
        })

        Context("with an invalid request", func() {
            It("returns false", func() {
                body := strings.NewReader(`{"text": "the text"}`)
                params := handlers.NewNotifySpaceParams(body)
                Expect(params.Validate()).To(BeFalse())
            })

            It("requires a kind attribute", func() {
                body := strings.NewReader(`{"text": "the text"}`)
                params := handlers.NewNotifySpaceParams(body)
                params.Validate()

                Expect(len(params.Errors)).To(Equal(1))
                Expect(params.Errors).To(ContainElement(`"kind" is a required field`))
            })

            It("requires a text attribute", func() {
                body := strings.NewReader(`{}`)
                params := handlers.NewNotifySpaceParams(body)
                params.Validate()

                Expect(len(params.Errors)).To(Equal(2))
                Expect(params.Errors).To(ContainElement(`"kind" is a required field`))
                Expect(params.Errors).To(ContainElement(`"text" is a required field`))
            })
        })
    })
})
