package handlers_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyParams", func() {
    Describe("NewNotifyParams", func() {
        It("parses the body of the given request", func() {
            body := strings.NewReader(`{
                "kind_id": "test_email",
                "reply_to": "me@awesome.com",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`)

            params, _ := handlers.NewNotifyParams(body)

            Expect(params.KindID).To(Equal("test_email"))
            Expect(params.KindDescription).To(Equal(""))
            Expect(params.SourceDescription).To(Equal(""))
            Expect(params.ReplyTo).To(Equal("me@awesome.com"))
            Expect(params.Subject).To(Equal("Summary of contents"))
            Expect(params.Text).To(Equal("Contents of the email message"))
        })

        It("does not blow up if the request body is empty", func() {
            body := strings.NewReader("")

            Expect(func() {
                handlers.NewNotifyParams(body)
            }).NotTo(Panic())
        })

        Describe("html parsing", func() {
            Context("when html is passed with a surrounding html and body tag", func() {
                It("pulls out the html in the body", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><head><title>BananaDamage</title></head><body><p>The TEXT</p><h1>the TITLE</h1></body></html>"
                    }`)

                    params, err := handlers.NewNotifyParams(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(params.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when html is passed with a surrounding body tag", func() {
                It("pulls out the html in the body", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<body><p>The TEXT</p><h1>the TITLE</h1></body>"
                    }`)

                    params, err := handlers.NewNotifyParams(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(params.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when html is passed with a surrounding html tag", func() {
                It("pulls out the html in the html tag", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT</p><h1>the TITLE</h1></html>"
                    }`)

                    params, err := handlers.NewNotifyParams(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(params.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when just bare html is passed without surrounding html/body tags", func() {
                It("is a no op", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<p>The TEXT</p><h1>the TITLE</h1>"
                    }`)

                    params, err := handlers.NewNotifyParams(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(params.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when invalid html is passed", func() {
                It("pulls out the html anyway", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT<h1>the TITLE</h1></html>"
                    }`)
                    params, err := handlers.NewNotifyParams(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(params.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))

                    body = strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT<h1>the TITLE</h1></body>"
                    }`)
                    params, err = handlers.NewNotifyParams(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(params.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when no html is passed", func() {
                It("does not error", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "text": "not html yo"
                    }`)
                    params, err := handlers.NewNotifyParams(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(params.HTML).To(Equal(""))
                })
            })
        })
    })

    Describe("Validate", func() {
        It("validates the required parameters in the request body", func() {
            body := strings.NewReader(`{
                "kind_id": "test_email",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`)
            params, err := handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))

            params.KindID = ""

            Expect(params.Validate()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(1))
            Expect(params.Errors).To(ContainElement(`"kind_id" is a required field`))

            params.Text = ""

            Expect(params.Validate()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(2))
            Expect(params.Errors).To(ContainElement(`"kind_id" is a required field`))
            Expect(params.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

            params.KindID = "something"
            params.Text = "banana"

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))
        })

        It("either text or html must be set", func() {
            body := strings.NewReader(`{
                "kind_id": "test_email"
            }`)

            params, err := handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            Expect(params.Validate()).To(BeFalse())
            Expect(params.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

            body = strings.NewReader(`{
                "kind_id": "test_email",
                "text": "Contents of the email message"
            }`)

            params, err = handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))

            body = strings.NewReader(`{
                "kind_id": "test_email",
                "html": "<html><body><p>the html</p></body></html>"
            }`)

            params, err = handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))

            body = strings.NewReader(`{
                "kind_id": "test_email",
                "text": "Contents of the email message",
                "html": "<html><body><p>the html</p></body></html>"
            }`)

            params, err = handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))
        })

        It("validates the format of kind_id", func() {
            body := strings.NewReader(`{
                "kind_id": "A_valid.id-99",
                "text": "Contents of the email message"
            }`)

            params, err := handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            Expect(params.Validate()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))

            body = strings.NewReader(`{
                "kind_id": "an_invalid.id-00!",
                "text": "Contents of the email message"
            }`)

            params, err = handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            Expect(params.Validate()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(1))
            Expect(params.Errors).To(ContainElement(`"kind_id" is improperly formatted`))
        })
    })

    Describe("ToOptions", func() {
        It("converts itself to a postal.Options object", func() {
            body := strings.NewReader(`{
                "kind_id": "test_email",
                "reply_to": "me@awesome.com",
                "subject": "Summary of contents",
                "text": "Contents of the email message",
                "html": "<div>Some HTML</div>"
            }`)

            params, err := handlers.NewNotifyParams(body)
            if err != nil {
                panic(err)
            }

            client := models.Client{
                ID:          "client-id",
                Description: "Descriptive Component Name",
            }
            kind := models.Kind{
                ID:          "test_email",
                ClientID:    "client-id",
                Description: "Descriptive Kind Name",
            }
            options := params.ToOptions(client, kind)
            Expect(options).To(Equal(postal.Options{
                KindID:            "test_email",
                KindDescription:   "Descriptive Kind Name",
                SourceDescription: "Descriptive Component Name",
                ReplyTo:           "me@awesome.com",
                Subject:           "Summary of contents",
                Text:              "Contents of the email message",
                HTML:              "<div>Some HTML</div>",
            }))
        })
    })
})
