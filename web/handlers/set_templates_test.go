package handlers_test

import (
    "bytes"
    "fmt"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("SetTemplates", func() {
    var err error
    var handler handlers.SetTemplates
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var context stack.Context
    var updater *fakes.FakeTemplateUpdater
    var fakeErrorWriter *fakes.FakeErrorWriter

    Describe("ServeHTTP", func() {
        BeforeEach(func() {
            updater = fakes.NewFakeTemplateUpdater()
            fakeErrorWriter = fakes.NewFakeErrorWriter()
            handler = handlers.NewSetTemplates(updater, fakeErrorWriter)
            writer = httptest.NewRecorder()
            body := []byte(`{"text": "{{turkey}}", "html": "<p>{{turkey}} gobble</p>"}`)
            request, err = http.NewRequest("PUT", "/templates/myTemplateName.user_body", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }
        })

        It("calls set on its setter with appropriate arguments", func() {
            handler.ServeHTTP(writer, request, context)
            Expect(updater.UpdateArgument).To(Equal(models.Template{
                Name:       "myTemplateName.user_body",
                Text:       "{{turkey}}",
                HTML:       "<p>{{turkey}} gobble</p>",
                Overridden: true,
            }))
            Expect(writer.Code).To(Equal(http.StatusNoContent))
        })

        Context("when an errors occurs", func() {
            It("Writes a 422 when the request is missing the text field", func() {
                body := []byte(`{"html": "<p>gobble</p>"}`)
                request, err = http.NewRequest("PUT", "/templates/myTemplateName.user_body", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(422))
            })

            It("Writes a 422 when the request is missing the html field", func() {
                body := []byte(`{"text": "gobble"}`)
                request, err = http.NewRequest("PUT", "/templates/myTemplateName.user_body", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(422))
            })

            It("writes a 400 for an invalid request", func() {
                body := []byte(`{"text": forgot to close the curly brace`)
                request, err = http.NewRequest("PUT", "/templates/myTemplateName.user_body", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(400))
            })

            It("returns a 500 for all other error cases", func() {
                updater.UpdateError = fmt.Errorf("my new error")
                handler.ServeHTTP(writer, request, context)
                Expect(writer.Code).To(Equal(http.StatusInternalServerError))
            })
        })
    })
})
