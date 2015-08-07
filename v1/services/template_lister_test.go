package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateLister", func() {
	var (
		lister            services.TemplateLister
		templatesRepo     *fakes.TemplatesRepo
		expectedTemplates map[string]services.TemplateSummary
		database          *fakes.Database
	)

	BeforeEach(func() {
		templatesRepo = fakes.NewTemplatesRepo()
		database = fakes.NewDatabase()

		lister = services.NewTemplateLister(templatesRepo)
	})

	Describe("List", func() {
		Context("when the templates exists in the database", func() {
			BeforeEach(func() {
				testTemplates := []models.Template{
					{
						ID:      "starwarr-guid",
						Name:    "Star Wars",
						Subject: "Awesomeness",
						HTML:    "<p>Millenium Falcon</p>",
						Text:    "Millenium Falcon",
					},
					{
						ID:      models.DefaultTemplateID,
						Name:    "default name",
						Subject: "default subject",
						HTML:    "<h1>default</h1>",
						Text:    "defaults!",
					},
					{
						ID:      "robot-guid",
						Name:    "Big Hero 6",
						Subject: "Heroes",
						HTML:    "<h1>Robots!</h1>",
						Text:    "Robots!",
					},
					{
						ID:      "boring-guid",
						Name:    "Blah",
						Subject: "More Blah",
						HTML:    "<h1>This is blahblah</h1>",
						Text:    "Blah even more",
					},
					{
						ID:      "starvation-guid",
						Name:    "Hungry Play",
						Subject: "Dystopian",
						HTML:    "<h1>Sad</h1>",
						Text:    "Run!!",
					},
				}
				expectedTemplates = map[string]services.TemplateSummary{
					"starwarr-guid":   {Name: "Star Wars"},
					"robot-guid":      {Name: "Big Hero 6"},
					"boring-guid":     {Name: "Blah"},
					"starvation-guid": {Name: "Hungry Play"},
				}
				templatesRepo.TemplatesList = testTemplates
			})

			It("returns a list of guids and template names", func() {
				templates, err := lister.List(database)
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(expectedTemplates))

				Expect(database.ConnectionWasCalled).To(BeTrue())
			})
		})

		Context("the lister has an error", func() {
			It("propagates the error", func() {
				templatesRepo.ListError = errors.New("some-error")
				_, err := lister.List(database)
				Expect(err.Error()).To(Equal("some-error"))
			})
		})
	})
})
