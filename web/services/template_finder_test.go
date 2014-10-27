package services_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"
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
            env.RootPath + "/templates/space_body.html":  "default-space-html",
            env.RootPath + "/templates/subject.missing":  "default-missing-subject",
            env.RootPath + "/templates/subject.provided": "default-provided-subject",
            env.RootPath + "/templates/user_body.text":   "default-user-text",
            env.RootPath + "/templates/user_body.html":   "default-user-html",
            env.RootPath + "/templates/email_body.html":  "email-body-html",
            env.RootPath + "/templates/email_body.text":  "email-body-text",
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

var _ = Describe("Finder", func() {
    var finder services.TemplateFinder
    var fakeTemplatesRepo *fakes.FakeTemplatesRepo
    var fakeFileSystem FakeFileSystem

    Describe("#Find", func() {
        BeforeEach(func() {
            env := config.NewEnvironment()
            fakeTemplatesRepo = fakes.NewFakeTemplatesRepo()
            fakeFileSystem = NewFakeFileSystem(env)

            finder = services.NewTemplateFinder(fakeTemplatesRepo, env.RootPath, fakes.NewDatabase(), fakeFileSystem)
        })

        Context("when the finder returns a template", func() {
            Context("when the override does not exist", func() {
                It("returns the default space template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp.space_body")
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("default-space-text"))
                    Expect(template.HTML).To(Equal("default-space-html"))
                })

                It("returns the default user template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp.user_body")
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("default-user-text"))
                    Expect(template.HTML).To(Equal("default-user-html"))
                })

                It("returns the default email template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp.email_body")
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("email-body-text"))
                    Expect(template.HTML).To(Equal("email-body-html"))
                })
            })

            Context("when the override exists in the database", func() {
                var expectedTemplate models.Template

                BeforeEach(func() {
                    expectedTemplate = models.Template{
                        Text:       "authenticate new hungry raptors template",
                        HTML:       "<p>hungry raptors are newly authenticated template</p>",
                        Overridden: true,
                    }
                    fakeTemplatesRepo.Templates["authentication.new.user_body"] = expectedTemplate
                })

                It("returns the requested override template", func() {
                    template, err := finder.Find("authentication.new.user_body")
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeTrue())
                    Expect(template).To(Equal(expectedTemplate))
                })

            })

            Context("when the requested client/kind override does not exist in db", func() {
                Context("but the client override does", func() {
                    var expectedTemplate models.Template

                    BeforeEach(func() {
                        expectedTemplate = models.Template{
                            Text:       "authentication template for hungry raptors",
                            HTML:       "<h1>Wow you are authentic!</h1>",
                            Overridden: true,
                        }
                        fakeTemplatesRepo.Templates["authentication.user_body"] = expectedTemplate
                    })

                    It("returns the fallback override that exists", func() {
                        template, err := finder.Find("authentication.new.user_body")
                        Expect(err).ToNot(HaveOccurred())
                        Expect(template.Overridden).To(BeTrue())
                        Expect(template).To(Equal(expectedTemplate))
                    })
                })

                Context("when the client override does not exist", func() {
                    var expectedTemplate models.Template

                    BeforeEach(func() {
                        expectedTemplate = models.Template{
                            Text:       "special user template",
                            HTML:       "<h1>Wow you are a special user!</h1>",
                            Overridden: true,
                        }
                        fakeTemplatesRepo.Templates["user_body"] = expectedTemplate
                    })

                    It("returns the fallback override that exists", func() {
                        template, err := finder.Find("authentication.new.user_body")
                        Expect(err).ToNot(HaveOccurred())
                        Expect(template.Overridden).To(BeTrue())
                        Expect(template).To(Equal(expectedTemplate))
                    })
                })
            })
        })

        Context("when the finder returns an error", func() {
            It("propagates the error", func() {
                fakeTemplatesRepo.FindError = errors.New("some-error")
                _, err := finder.Find("missing_template_file")
                Expect(err.Error()).To(Equal("some-error"))
            })
        })
    })
})
