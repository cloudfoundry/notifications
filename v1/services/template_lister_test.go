package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateLister", func() {
	var (
		lister        services.TemplateLister
		templatesRepo *mocks.TemplatesRepo
		database      *mocks.Database
		conn          *mocks.Connection
	)

	BeforeEach(func() {
		templatesRepo = mocks.NewTemplatesRepo()
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		lister = services.NewTemplateLister(templatesRepo)
	})

	Describe("List", func() {
		Context("when the templates exists in the database", func() {
			BeforeEach(func() {
				templatesRepo.ListIDsAndNamesCall.Returns.Templates = []models.Template{
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
			})

			It("returns a list of guids and template names", func() {
				templates, err := lister.List(database)
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(map[string]services.TemplateSummary{
					"starwarr-guid":   {Name: "Star Wars"},
					"robot-guid":      {Name: "Big Hero 6"},
					"boring-guid":     {Name: "Blah"},
					"starvation-guid": {Name: "Hungry Play"},
				}))

				Expect(templatesRepo.ListIDsAndNamesCall.Receives.Connection).To(Equal(conn))
			})
		})

		Context("the lister has an error", func() {
			It("propagates the error", func() {
				templatesRepo.ListIDsAndNamesCall.Returns.Error = errors.New("some-error")

				_, err := lister.List(database)
				Expect(err).To(MatchError(errors.New("some-error")))
			})
		})
	})
})
