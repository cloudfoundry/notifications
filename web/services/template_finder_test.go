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

var _ = Describe("Finder", func() {
	var finder services.TemplateFinder
	var templatesRepo *fakes.TemplatesRepo
	var fileSystem fakes.FileSystem

	Describe("#FindByID", func() {
		BeforeEach(func() {
			env := config.NewEnvironment()
			templatesRepo = fakes.NewTemplatesRepo()
			fileSystem = fakes.NewFileSystem(env.RootPath)

			finder = services.NewTemplateFinder(templatesRepo, env.RootPath, fakes.NewDatabase(), fileSystem)
		})

		Context("when the finder returns a template", func() {
			Context("when the template exists in the database", func() {
				var expectedTemplate models.Template

				BeforeEach(func() {
					expectedTemplate = models.Template{
						ID:      "awesome-template-id",
						Name:    "Awesome New Template",
						Subject: "Wow this is really awesome",
						Text:    "awesome new hungry raptors template",
						HTML:    "<p>hungry raptors are newly awesome template</p>",
					}
					templatesRepo.Templates["awesome-template-id"] = expectedTemplate
				})

				It("returns the requested template", func() {
					template, err := finder.FindByID("awesome-template-id")
					Expect(err).ToNot(HaveOccurred())
					Expect(template).To(Equal(expectedTemplate))
				})
			})

		})

		Context("the finder has an error", func() {
			It("propagates the error", func() {
				templatesRepo.FindError = errors.New("some-error")
				templatesRepo.Templates["some-template-id"] = models.Template{
					Name:    "Not nice template",
					Subject: "Not the kind you want",
					Text:    "throwing errors template",
					HTML:    "<h1>Wow you are a throwing errors!</h1>",
				}
				_, err := finder.FindByID("some-template-id")
				Expect(err.Error()).To(Equal("some-error"))
			})
		})
	})

	Describe("#Find", func() {
		BeforeEach(func() {
			env := config.NewEnvironment()
			templatesRepo = fakes.NewTemplatesRepo()
			fileSystem = fakes.NewFileSystem(env.RootPath)

			finder = services.NewTemplateFinder(templatesRepo, env.RootPath, fakes.NewDatabase(), fileSystem)
		})

		Context("when the finder returns a template", func() {
			Context("when the override does not exist", func() {
				It("returns the default space template", func() {
					templatesRepo.FindError = models.RecordNotFoundError("")

					template, err := finder.Find("login.fp." + models.SpaceBodyTemplateName)
					Expect(err).ToNot(HaveOccurred())
					Expect(template.Text).To(Equal("default-space-text"))
					Expect(template.HTML).To(Equal("default-space-html"))
				})

				It("returns the default user template", func() {
					templatesRepo.FindError = models.RecordNotFoundError("")

					template, err := finder.Find("login.fp." + models.UserBodyTemplateName)
					Expect(err).ToNot(HaveOccurred())
					Expect(template.Text).To(Equal("default-user-text"))
					Expect(template.HTML).To(Equal("default-user-html"))
				})

				It("returns the default email template", func() {
					templatesRepo.FindError = models.RecordNotFoundError("")

					template, err := finder.Find("login.fp." + models.EmailBodyTemplateName)
					Expect(err).ToNot(HaveOccurred())
					Expect(template.Text).To(Equal("email-body-text"))
					Expect(template.HTML).To(Equal("email-body-html"))
				})

				It("returns the default subject missing template", func() {
					templatesRepo.FindError = models.RecordNotFoundError("")

					template, err := finder.Find("login.fp." + models.SubjectMissingTemplateName)
					Expect(err).ToNot(HaveOccurred())
					Expect(template.Text).To(Equal("default-missing-subject"))

				})

				It("returns the default subject provided template", func() {
					templatesRepo.FindError = models.RecordNotFoundError("")

					template, err := finder.Find("login.fp." + models.SubjectProvidedTemplateName)
					Expect(err).ToNot(HaveOccurred())
					Expect(template.Text).To(Equal("default-provided-subject"))

				})
			})

			Context("when the override exists in the database", func() {
				var expectedTemplate models.Template

				BeforeEach(func() {
					expectedTemplate = models.Template{
						Name: "authentication.new." + models.UserBodyTemplateName,
						Text: "authenticate new hungry raptors template",
						HTML: "<p>hungry raptors are newly authenticated template</p>",
					}
					templatesRepo.Templates["authentication.new."+models.UserBodyTemplateName] = expectedTemplate
				})

				It("returns the requested override template", func() {
					template, err := finder.Find("authentication.new." + models.UserBodyTemplateName)
					Expect(err).ToNot(HaveOccurred())
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
						templatesRepo.Templates["authentication."+models.UserBodyTemplateName] = expectedTemplate
					})

					It("returns the fallback override that exists", func() {
						template, err := finder.Find("authentication.new." + models.UserBodyTemplateName)
						Expect(err).ToNot(HaveOccurred())
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
						templatesRepo.Templates[models.UserBodyTemplateName] = expectedTemplate
					})

					It("returns the fallback override that exists", func() {
						template, err := finder.Find("authentication.new." + models.UserBodyTemplateName)
						Expect(err).ToNot(HaveOccurred())
						Expect(template).To(Equal(expectedTemplate))
					})
				})
			})
		})

		Context("the finder has an error", func() {
			It("propagates the error", func() {
				templatesRepo.FindError = errors.New("some-error")
				templatesRepo.Templates[models.UserBodyTemplateName] = models.Template{
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
				"login.fp." + models.UserBodyTemplateName:        []string{"login.fp." + models.UserBodyTemplateName, "login." + models.UserBodyTemplateName, models.UserBodyTemplateName},
				"login." + models.UserBodyTemplateName:           []string{"login." + models.UserBodyTemplateName, models.UserBodyTemplateName},
				models.UserBodyTemplateName:                      []string{models.UserBodyTemplateName},
				"login.fp." + models.SubjectMissingTemplateName:  []string{"login.fp." + models.SubjectMissingTemplateName, "login." + models.SubjectMissingTemplateName, models.SubjectMissingTemplateName},
				"login." + models.SubjectMissingTemplateName:     []string{"login." + models.SubjectMissingTemplateName, models.SubjectMissingTemplateName},
				models.SubjectMissingTemplateName:                []string{models.SubjectMissingTemplateName},
				"login.fp.banana." + models.UserBodyTemplateName: []string{},
				"banana": []string{},
			}

			for input, output := range table {
				names := finder.ParseTemplateName(input)
				Expect(names).To(Equal(output))
			}
		})
	})
})
