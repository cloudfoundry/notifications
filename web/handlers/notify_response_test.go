package handlers_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyResponseGenerator", func() {
    Describe("LoadTemplates", func() {
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

                plain, html, err := notifyResponse.LoadTemplates(true, manager)
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

                plain, html, err := notifyResponse.LoadTemplates(false, manager)
                if err != nil {
                    panic(err)
                }

                Expect(plain).To(Equal("user plain text"))
                Expect(html).To(Equal("user html code"))
            })
        })
    })
})
