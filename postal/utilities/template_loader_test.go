package utilities_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal/utilities"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("TemplateLoader", func() {
    var loader utilities.TemplatesLoader
    var finder *fakes.TemplateFinder

    BeforeEach(func() {
        finder = fakes.NewTemplateFinder()

        finder.Templates["raptors.hungry."+models.SubjectProvidedTemplateName] = models.Template{
            Text: "Dinosaurs are coming",
        }

        finder.Templates["raptors.hungry."+models.UserBodyTemplateName] = models.Template{
            HTML: "<p>Can Raptors Open Doors?</p>",
            Text: "Yes they ca--",
        }

        loader = utilities.NewTemplatesLoader(finder)
    })

    Describe("LoadTemplates", func() {

        It("Returns templates using its finder", func() {
            templates, err := loader.LoadTemplates(models.SubjectProvidedTemplateName, models.UserBodyTemplateName, "raptors", "hungry")
            Expect(err).ToNot(HaveOccurred())
            Expect(templates.HTML).To(Equal("<p>Can Raptors Open Doors?</p>"))
            Expect(templates.Text).To(Equal("Yes they ca--"))
            Expect(templates.Subject).To(Equal("Dinosaurs are coming"))
        })

        Context("The finder errors", func() {
            It("Propagates that error", func() {
                finder.FindError = errors.New("Boom!")
                _, err := loader.LoadTemplates(models.SubjectProvidedTemplateName, models.UserBodyTemplateName, "raptors", "hungry")
                Expect(err).To(HaveOccurred())

            })
        })

    })

})
