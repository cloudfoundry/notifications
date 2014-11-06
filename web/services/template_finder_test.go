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
            env.RootPath + "/templates/" + models.SpaceBodyTemplateName + ".text": "default-space-text",
            env.RootPath + "/templates/" + models.SpaceBodyTemplateName + ".html": "default-space-html",
            env.RootPath + "/templates/" + models.SubjectMissingTemplateName:      "default-missing-subject",
            env.RootPath + "/templates/" + models.SubjectProvidedTemplateName:     "default-provided-subject",
            env.RootPath + "/templates/" + models.UserBodyTemplateName + ".text":  "default-user-text",
            env.RootPath + "/templates/" + models.UserBodyTemplateName + ".html":  "default-user-html",
            env.RootPath + "/templates/" + models.EmailBodyTemplateName + ".html": "email-body-html",
            env.RootPath + "/templates/" + models.EmailBodyTemplateName + ".text": "email-body-text",
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

                    template, err := finder.Find("login.fp." + models.SpaceBodyTemplateName)
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("default-space-text"))
                    Expect(template.HTML).To(Equal("default-space-html"))
                })

                It("returns the default user template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp." + models.UserBodyTemplateName)
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("default-user-text"))
                    Expect(template.HTML).To(Equal("default-user-html"))
                })

                It("returns the default email template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp." + models.EmailBodyTemplateName)
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("email-body-text"))
                    Expect(template.HTML).To(Equal("email-body-html"))
                })

                It("returns the default subject missing template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp." + models.SubjectMissingTemplateName)
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("default-missing-subject"))

                })

                It("returns the default subject provided template", func() {
                    fakeTemplatesRepo.FindError = models.ErrRecordNotFound{}

                    template, err := finder.Find("login.fp." + models.SubjectProvidedTemplateName)
                    Expect(err).ToNot(HaveOccurred())
                    Expect(template.Overridden).To(BeFalse())
                    Expect(template.Text).To(Equal("default-provided-subject"))

                })
            })

            Context("when the override exists in the database", func() {
                var expectedTemplate models.Template

                BeforeEach(func() {
                    expectedTemplate = models.Template{
                        Name:       "authentication.new." + models.UserBodyTemplateName,
                        Text:       "authenticate new hungry raptors template",
                        HTML:       "<p>hungry raptors are newly authenticated template</p>",
                        Overridden: true,
                    }
                    fakeTemplatesRepo.Templates["authentication.new."+models.UserBodyTemplateName] = expectedTemplate
                })

                It("returns the requested override template", func() {
                    template, err := finder.Find("authentication.new." + models.UserBodyTemplateName)
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
                            Name: "authentication." + models.UserBodyTemplateName,
                            Text: "authentication template for hungry raptors",
                            HTML: "<h1>Wow you are authentic!</h1>",
                        }
                        fakeTemplatesRepo.Templates["authentication."+models.UserBodyTemplateName] = expectedTemplate
                    })

                    It("returns the fallback override that exists", func() {
                        template, err := finder.Find("authentication.new." + models.UserBodyTemplateName)
                        Expect(err).ToNot(HaveOccurred())
                        Expect(template.Overridden).To(BeFalse())
                        Expect(template).To(Equal(expectedTemplate))
                    })
                })

                Context("when the client override does not exist", func() {
                    var expectedTemplate models.Template

                    BeforeEach(func() {
                        expectedTemplate = models.Template{
                            Name: models.UserBodyTemplateName,
                            Text: "special user template",
                            HTML: "<h1>Wow you are a special user!</h1>",
                        }
                        fakeTemplatesRepo.Templates[models.UserBodyTemplateName] = expectedTemplate
                    })

                    It("returns the fallback override that exists", func() {
                        template, err := finder.Find("authentication.new." + models.UserBodyTemplateName)
                        Expect(err).ToNot(HaveOccurred())
                        Expect(template.Overridden).To(BeFalse())
                        Expect(template).To(Equal(expectedTemplate))
                    })
                })
            })
        })

        Context("the finder has an error", func() {
            It("propagates the error", func() {
                fakeTemplatesRepo.FindError = errors.New("some-error")
                fakeTemplatesRepo.Templates[models.UserBodyTemplateName] = models.Template{
                    Name: models.UserBodyTemplateName,
                    Text: "special user template",
                    HTML: "<h1>Wow you are a special user!</h1>",
                }
                _, err := finder.Find(models.UserBodyTemplateName)
                Expect(err.Error()).To(Equal("some-error"))
            })
        })

        Context("when the template name does not match a known format", func() {
            It("returns a TemplateNotFound error", func() {
                _, err := finder.Find("banana")
                Expect(err).To(BeAssignableToTypeOf(services.TemplateNotFoundError("")))
            })
        })
    })

    Describe("#ParseTemplateName", func() {
        It("parses the input template name, returning a list of possible template matches", func() {
            table := map[string][]string{
                "login.fp.user_body": []string{"login.fp.user_body", "login.user_body", "user_body"},
                "login.user_body": []string{"login.user_body", "user_body"},
                "user_body": []string{"user_body"},
                "login.fp.subject.missing": []string{"login.fp.subject.missing", "login.subject.missing", "subject.missing"},
                "login.subject.missing": []string{"login.subject.missing", "subject.missing"},
                "subject.missing": []string{"subject.missing"},
                "banana": []string{},
                "login.fp.banana.user_body": []string{},
            }

            for input, output := range table {
                names := finder.ParseTemplateName(input)
                Expect(names).To(Equal(output))
            }
        })
    })
})
