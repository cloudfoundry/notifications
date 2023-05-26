package v1_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/postal/v1"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateLoader", func() {
	var (
		loader        v1.TemplatesLoader
		clientsRepo   *mocks.ClientsRepository
		kindsRepo     *mocks.KindsRepo
		templatesRepo *mocks.TemplatesRepo
		conn          db.ConnectionInterface
		database      *mocks.Database
	)

	BeforeEach(func() {
		clientsRepo = mocks.NewClientsRepository()
		kindsRepo = mocks.NewKindsRepo()
		templatesRepo = mocks.NewTemplatesRepo()

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		loader = v1.NewTemplatesLoader(database, clientsRepo, kindsRepo, templatesRepo)
	})

	Describe("LoadTemplates", func() {
		BeforeEach(func() {
			clientsRepo.FindCall.Returns.Client = models.Client{
				ID:         "my-client-id",
				TemplateID: models.DefaultTemplateID,
			}

			kindsRepo.FindCall.Returns.Kinds = []models.Kind{
				{
					ID:         "my-kind-id",
					ClientID:   "my-client-id",
					TemplateID: models.DefaultTemplateID,
				},
			}

			templatesRepo.FindByIDCall.Returns.Template = models.Template{
				ID:      models.DefaultTemplateID,
				Name:    "Default Template",
				HTML:    "<p>The default template</p>",
				Text:    "The default template",
				Subject: "default subject",
			}
		})

		Context("when the kind has a template", func() {
			BeforeEach(func() {
				templatesRepo.FindByIDCall.Returns.Template = models.Template{
					ID:      "my-kind-template",
					Name:    "my-kind-template",
					HTML:    "<p>kind template</p>",
					Text:    "some kind template text",
					Subject: "kind subject",
				}

				kindsRepo.FindCall.Returns.Kinds = []models.Kind{
					{
						ID:         "my-kind-id",
						ClientID:   "my-client-id",
						TemplateID: "my-kind-template",
					},
				}
			})

			It("returns the template belonging to the kind", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(common.Templates{
					HTML:    "<p>kind template</p>",
					Text:    "some kind template text",
					Subject: "kind subject",
				}))

				Expect(templatesRepo.FindByIDCall.Receives.Connection).To(Equal(conn))
				Expect(templatesRepo.FindByIDCall.Receives.TemplateID).To(Equal("my-kind-template"))
			})
		})

		Context("when the client has a template", func() {
			BeforeEach(func() {
				templatesRepo.FindByIDCall.Returns.Template = models.Template{
					ID:      "my-client-template",
					Name:    "my-client-template",
					HTML:    "<p>client template</p>",
					Text:    "some client template text",
					Subject: "client subject",
				}

				clientsRepo.FindCall.Returns.Client = models.Client{
					ID:         "my-client-id",
					TemplateID: "my-client-template",
				}
			})

			It("returns the template belonging to the client", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(common.Templates{
					HTML:    "<p>client template</p>",
					Text:    "some client template text",
					Subject: "client subject",
				}))

				Expect(templatesRepo.FindByIDCall.Receives.Connection).To(Equal(conn))
				Expect(templatesRepo.FindByIDCall.Receives.TemplateID).To(Equal("my-client-template"))
			})
		})

		Context("when the neither client nor kind has a template", func() {
			It("returns the default template", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(common.Templates{
					HTML:    "<p>The default template</p>",
					Text:    "The default template",
					Subject: "default subject",
				}))
			})
		})

		Context("when kindID is an empty string", func() {
			It("does not look for a template belonging to the kind", func() {
				templates, err := loader.LoadTemplates("my-client-id", "", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(common.Templates{
					HTML:    "<p>The default template</p>",
					Text:    "The default template",
					Subject: "default subject",
				}))
			})
		})

		Context("when the kinds repo has an error", func() {
			It("bubbles up the error", func() {
				kindsRepo.FindCall.Returns.Error = errors.New("BOOM!")

				_, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).To(HaveOccurred())
			})

		})

		Context("when the clients repo has an error", func() {
			It("bubbles up the error", func() {
				clientsRepo.FindCall.Returns.Error = errors.New("BOOM!")

				_, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
