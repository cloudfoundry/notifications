package handlers_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("EmailTemplateManager", func() {
    var manager handlers.EmailTemplateManager

    Describe("LoadEmailTemplate", func() {
        Context("when there are no template overrides", func() {
            It("loads the templates from the default location", func() {
                manager.ReadFile = func(path string) (string, error) {
                    if strings.Contains(path, "text") {
                        return "the fake text", nil
                    }
                    return "incorrect", nil
                }

                manager.FileExists = func(path string) bool {
                    return false
                }

                plainTextTemplate, err := manager.LoadEmailTemplate("user_body.text")
                if err != nil {
                    panic(err)
                }
                Expect(plainTextTemplate).To(Equal("the fake text"))
            })
        })

        Context("when a template has an override set", func() {
            It("replaces the default template with the user provided one", func() {

                manager.ReadFile = func(path string) (string, error) {
                    if strings.Contains(path, "overrides") {
                        return "the fake text", nil
                    }
                    return "incorrect", nil
                }

                manager.FileExists = func(path string) bool {
                    return true
                }

                userTextEmail, err := manager.LoadEmailTemplate("user_body.text")
                if err != nil {
                    panic(err)
                }

                Expect(userTextEmail).To(Equal("the fake text"))
            })
        })
    })
})
