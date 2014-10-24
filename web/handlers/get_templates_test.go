package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/ryanmoran/stack"
)

var _ = Describe("GetTemplates", func() {
    var handler handlers.GetTemplates
    var request *http.Request
    var writer *httptest.ResponseRecorder
    var context stack.Context
    var finder *fakes.FakeTemplateFinder

    Describe("ServeHTTP", func() {

        Context("With a proper template name", func() {
            BeforeEach(func() {
                finder = fakes.NewFakeTemplateFinder(models.Template{
                    Text: "the template {{variable}}",
                    HTML: "<p> the template {{variable}} </p>",
                })

                writer = httptest.NewRecorder()
                handler = handlers.NewGetTemplates(finder)

                var err error
                request, err = http.NewRequest("GET", "/templates/myTemplateName.user_body", bytes.NewBuffer([]byte{}))
                if err != nil {
                    panic(err)
                }
            })

            It("Calls Execute on its finder with appropriate arguments", func() {
                handler.ServeHTTP(writer, request, context)
                Expect(finder.TemplateName).To(Equal("myTemplateName.user_body"))
            })

            It("writes out the finder's response", func() {
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(http.StatusOK))

                var template models.Template
                err := json.Unmarshal(writer.Body.Bytes(), &template)
                if err != nil {
                    panic(err)
                }

                Expect(template.HTML).To(Equal("<p> the template {{variable}} </p>"))
                Expect(template.Text).To(Equal("the template {{variable}}"))
            })
        })

        Context("With improper template name", func() {
            BeforeEach(func() {
                writer = httptest.NewRecorder()
                handler = handlers.NewGetTemplates(finder)

                finder.FindError = models.ErrRecordNotFound{}
                var err error
                request, err = http.NewRequest("GET", "/templates/myTemplateName", bytes.NewBuffer([]byte{}))
                if err != nil {
                    panic(err)
                }
            })

            It("returns a 404 when the file name does not end with user_body or space_body", func() {
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(http.StatusNotFound))
            })
        })

        Context("When the finder errors", func() {
            BeforeEach(func() {
                writer = httptest.NewRecorder()
                finder = fakes.NewFakeTemplateFinder(models.Template{
                    Text: "the template {{variable}}",
                    HTML: "<p> the template {{variable}} </p>",
                })
                handler = handlers.NewGetTemplates(finder)
                finder.FindError = errors.New("BOOM!")

                var err error
                request, err = http.NewRequest("GET", "/templates/myTemplateName.user_body", bytes.NewBuffer([]byte{}))
                if err != nil {
                    panic(err)
                }
            })

            It("returns a 500", func() {
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(http.StatusInternalServerError))

            })
        })
    })
})
