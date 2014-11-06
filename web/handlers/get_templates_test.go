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
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetTemplates", func() {
    var handler handlers.GetTemplates
    var request *http.Request
    var writer *httptest.ResponseRecorder
    var context stack.Context
    var finder *fakes.FakeTemplateFinder
    var fakeErrorWriter *fakes.FakeErrorWriter

    Describe("ServeHTTP", func() {
        BeforeEach(func() {
            finder = fakes.NewFakeTemplateFinder()

            finder.Templates["myTemplateName."+models.UserBodyTemplateName] = models.Template{
                Text: "the template {{variable}}",
                HTML: "<p> the template {{variable}} </p>",
            }
            writer = httptest.NewRecorder()
            fakeErrorWriter = fakes.NewFakeErrorWriter()
            handler = handlers.NewGetTemplates(finder, fakeErrorWriter)
        })

        Context("When the finder does not error", func() {
            BeforeEach(func() {

                var err error
                request, err = http.NewRequest("GET", "/templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer([]byte{}))
                if err != nil {
                    panic(err)
                }
            })

            It("Calls find on its finder with appropriate arguments", func() {
                handler.ServeHTTP(writer, request, context)
                Expect(finder.TemplateName).To(Equal("myTemplateName." + models.UserBodyTemplateName))
            })

            It("writes out the finder's response", func() {
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(http.StatusOK))

                var template map[string]interface{}

                err := json.Unmarshal(writer.Body.Bytes(), &template)
                if err != nil {
                    panic(err)
                }

                Expect(len(template)).To(Equal(3))
                Expect(template["html"].(string)).To(Equal("<p> the template {{variable}} </p>"))
                Expect(template["text"].(string)).To(Equal("the template {{variable}}"))
                Expect(template["overridden"].(bool)).To(Equal(false))
            })
        })

        Context("When the finder errors", func() {
            BeforeEach(func() {
                finder.FindError = errors.New("BOOM!")

                var err error
                request, err = http.NewRequest("GET", "/templates/myTemplateName."+models.UserBodyTemplateName, bytes.NewBuffer([]byte{}))
                if err != nil {
                    panic(err)
                }
            })

            It("writes the error to the errorWriter", func() {
                handler.ServeHTTP(writer, request, context)
                Expect(fakeErrorWriter.Error).To(Equal(errors.New("BOOM!")))
            })
        })
    })
})
