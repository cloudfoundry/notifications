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

	Describe("#Find", func() {
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
					template, err := finder.Find("awesome-template-id")
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
				_, err := finder.Find("some-template-id")
				Expect(err.Error()).To(Equal("some-error"))
			})
		})

		Context("when the template cannot be found in the database", func() {
			XIt("returns a TemplateNotFound error", func() {
				_, err := finder.Find("banana")
				Expect(err).To(BeAssignableToTypeOf(services.TemplateNotFoundError("")))
			})
		})
	})

})
