package handlers_test

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyResponseGenerator", func() {

    Describe("LoadSubjectTemplate", func() {
        var manager handlers.EmailTemplateManager
        var notifyResponse handlers.NotifyResponseGenerator

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
