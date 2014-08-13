package params_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers/params"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Notify", func() {
    Describe("NewNotify", func() {
        It("parses the body of the given request", func() {
            body := strings.NewReader(`{
                "kind_id": "test_email",
                "reply_to": "me@awesome.com",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`)

            parameters, _ := params.NewNotify(body)

            Expect(parameters.KindID).To(Equal("test_email"))
            Expect(parameters.KindDescription).To(Equal(""))
            Expect(parameters.SourceDescription).To(Equal(""))
            Expect(parameters.ReplyTo).To(Equal("me@awesome.com"))
            Expect(parameters.Subject).To(Equal("Summary of contents"))
            Expect(parameters.Text).To(Equal("Contents of the email message"))
        })

        It("does not blow up if the request body is empty", func() {
            body := strings.NewReader("")

            Expect(func() {
                params.NewNotify(body)
            }).NotTo(Panic())
        })

        Describe("html parsing", func() {
            Context("when html is passed with a surrounding html and body tag", func() {
                It("pulls out the html in the body", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><head><title>BananaDamage</title></head><body><p>The TEXT</p><h1>the TITLE</h1></body></html>"
                    }`)

                    parameters, err := params.NewNotify(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parameters.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when html is passed with a surrounding body tag", func() {
                It("pulls out the html in the body", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<body><p>The TEXT</p><h1>the TITLE</h1></body>"
                    }`)

                    parameters, err := params.NewNotify(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parameters.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when html is passed with a surrounding html tag", func() {
                It("pulls out the html in the html tag", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT</p><h1>the TITLE</h1></html>"
                    }`)

                    parameters, err := params.NewNotify(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parameters.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when just bare html is passed without surrounding html/body tags", func() {
                It("is a no op", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<p>The TEXT</p><h1>the TITLE</h1>"
                    }`)

                    parameters, err := params.NewNotify(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parameters.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when invalid html is passed", func() {
                It("pulls out the html anyway", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT<h1>the TITLE</h1></html>"
                    }`)
                    parameters, err := params.NewNotify(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parameters.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))

                    body = strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT<h1>the TITLE</h1></body>"
                    }`)
                    parameters, err = params.NewNotify(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parameters.HTML).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
                })
            })

            Context("when no html is passed", func() {
                It("does not error", func() {
                    body := strings.NewReader(`{
                        "kind_id": "test_email",
                        "text": "not html yo"
                    }`)
                    parameters, err := params.NewNotify(body)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parameters.HTML).To(Equal(""))
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
            parameters, err := params.NewNotify(body)
            if err != nil {
                panic(err)
            }

            Expect(parameters.Validate()).To(BeTrue())
            Expect(len(parameters.Errors)).To(Equal(0))

            parameters.KindID = ""

            Expect(parameters.Validate()).To(BeFalse())
            Expect(len(parameters.Errors)).To(Equal(1))
            Expect(parameters.Errors).To(ContainElement(`"kind_id" is a required field`))

            parameters.Text = ""

            Expect(parameters.Validate()).To(BeFalse())
            Expect(len(parameters.Errors)).To(Equal(2))
            Expect(parameters.Errors).To(ContainElement(`"kind_id" is a required field`))
            Expect(parameters.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

            parameters.KindID = "something"
            parameters.Text = "banana"

            Expect(parameters.Validate()).To(BeTrue())
            Expect(len(parameters.Errors)).To(Equal(0))
        })

        It("either text or html must be set", func() {
            body := strings.NewReader(`{
                "kind_id": "test_email"
            }`)

            parameters, err := params.NewNotify(body)
            if err != nil {
                panic(err)
            }

            Expect(parameters.Validate()).To(BeFalse())
            Expect(parameters.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

            body = strings.NewReader(`{
                "kind_id": "test_email",
                "text": "Contents of the email message"
            }`)

            parameters, err = params.NewNotify(body)
            if err != nil {
                panic(err)
            }

            Expect(parameters.Validate()).To(BeTrue())
            Expect(len(parameters.Errors)).To(Equal(0))

            body = strings.NewReader(`{
                "kind_id": "test_email",
                "html": "<html><body><p>the html</p></body></html>"
            }`)

            parameters, err = params.NewNotify(body)
            if err != nil {
                panic(err)
            }

            Expect(parameters.Validate()).To(BeTrue())
            Expect(len(parameters.Errors)).To(Equal(0))

            body = strings.NewReader(`{
                "kind_id": "test_email",
                "text": "Contents of the email message",
                "html": "<html><body><p>the html</p></body></html>"
            }`)

            parameters, err = params.NewNotify(body)
            if err != nil {
                panic(err)
            }

            Expect(parameters.Validate()).To(BeTrue())
            Expect(len(parameters.Errors)).To(Equal(0))
        })

        It("validates the format of kind_id", func() {
            body := strings.NewReader(`{
                "kind_id": "A_valid.id-99",
                "text": "Contents of the email message"
            }`)

            parameters, err := params.NewNotify(body)
            if err != nil {
                panic(err)
            }

            Expect(parameters.Validate()).To(BeTrue())
            Expect(len(parameters.Errors)).To(Equal(0))

            body = strings.NewReader(`{
                "kind_id": "an_invalid.id-00!",
                "text": "Contents of the email message"
            }`)

            parameters, err = params.NewNotify(body)
            if err != nil {
                panic(err)
            }

            Expect(parameters.Validate()).To(BeFalse())
            Expect(len(parameters.Errors)).To(Equal(1))
            Expect(parameters.Errors).To(ContainElement(`"kind_id" is improperly formatted`))
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

            parameters, err := params.NewNotify(body)
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
            options := parameters.ToOptions(client, kind)
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
