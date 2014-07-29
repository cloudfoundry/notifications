package postal_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/postal"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

type FakeFileSystem struct {
    Files map[string]string
}

func NewFakeFileSystem(env config.Environment) FakeFileSystem {
    return FakeFileSystem{
        Files: map[string]string{
            env.RootPath + "/templates/space_body.text":  "default-space-text",
            env.RootPath + "/templates/subject.missing":  "default-missing-subject",
            env.RootPath + "/templates/space_body.html":  "default-space-html",
            env.RootPath + "/templates/subject.provided": "default-provided-subject",
            env.RootPath + "/templates/user_body.text":   "default-user-text",
            env.RootPath + "/templates/user_body.html":   "default-user-html",
        },
    }
}

func (fs FakeFileSystem) Exists(path string) bool {
    _, ok := fs.Files[path]
    return ok
}

func (fs FakeFileSystem) Read(path string) (string, error) {
    if file, ok := fs.Files[path]; ok {
        return file, nil
    }
    return "", errors.New("File does not exist")
}

var _ = Describe("TemplateLoader", func() {
    var loader postal.TemplateLoader
    var fs FakeFileSystem
    var env config.Environment

    BeforeEach(func() {
        env = config.NewEnvironment()
        fs = NewFakeFileSystem(env)
        loader = postal.NewTemplateLoader(&fs)
    })

    Describe("Load", func() {
        Context("when subject is not set in the params", func() {
            It("returns the subject.missing template", func() {
                templates, err := loader.Load("", postal.IsSpace)
                if err != nil {
                    panic(err)
                }

                Expect(templates.Subject).To(Equal("default-missing-subject"))
            })
        })

        Context("when subject is set in the params", func() {
            It("returns the subject.provided template", func() {
                templates, err := loader.Load("is provided", postal.IsSpace)
                if err != nil {
                    panic(err)
                }

                Expect(templates.Subject).To(Equal("default-provided-subject"))
            })
        })

        Context("notificationType is IsSpace", func() {
            It("returns the space templates", func() {
                templates, err := loader.Load("", postal.IsSpace)
                if err != nil {
                    panic(err)
                }

                Expect(templates.Text).To(Equal("default-space-text"))
                Expect(templates.HTML).To(Equal("default-space-html"))
            })
        })

        Context("notificationType is IsSpace", func() {
            It("returns the user templates", func() {
                templates, err := loader.Load("", postal.IsUser)
                if err != nil {
                    panic(err)
                }

                Expect(templates.Text).To(Equal("default-user-text"))
                Expect(templates.HTML).To(Equal("default-user-html"))
            })
        })
    })

    Describe("LoadTemplate", func() {
        Context("when there are no template overrides", func() {
            It("loads the templates from the default location", func() {
                text, err := loader.LoadTemplate("user_body.text")
                if err != nil {
                    panic(err)
                }
                Expect(text).To(Equal("default-user-text"))
            })
        })

        Context("when a template has an override set", func() {
            BeforeEach(func() {
                fs.Files[env.RootPath+"/templates/overrides/user_body.text"] = "override-user-text"
            })

            It("replaces the default template with the user provided one", func() {
                text, err := loader.LoadTemplate("user_body.text")
                if err != nil {
                    panic(err)
                }

                Expect(text).To(Equal("override-user-text"))
            })
        })
    })
})
