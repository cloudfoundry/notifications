package postal_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateLoader", func() {
	var (
		loader              postal.TemplatesLoader
		clientsRepo         *mocks.ClientsRepository
		kindsRepo           *mocks.KindsRepo
		templatesRepo       *mocks.TemplatesRepo
		conn                db.ConnectionInterface
		database            *mocks.Database
		templatesCollection *mocks.TemplatesCollection
	)

	BeforeEach(func() {
		clientsRepo = mocks.NewClientsRepository()
		kindsRepo = mocks.NewKindsRepo()
		templatesRepo = mocks.NewTemplatesRepo()
		database = mocks.NewDatabase()
		conn = database.Connection()
		templatesCollection = mocks.NewTemplatesCollection()
		loader = postal.NewTemplatesLoader(database, clientsRepo, kindsRepo, templatesRepo, templatesCollection)
	})

	Describe("LoadTemplates", func() {
		var kind models.Kind

		BeforeEach(func() {
			var err error

			clientsRepo.FindCall.Returns.Client = models.Client{
				ID:         "my-client-id",
				TemplateID: models.DefaultTemplateID,
			}

			kind, err = kindsRepo.Create(conn, models.Kind{
				ID:       "my-kind-id",
				ClientID: "my-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = templatesRepo.Create(conn, models.Template{
				ID:      models.DefaultTemplateID,
				Name:    "Default Template",
				HTML:    "<p>The default template</p>",
				Text:    "The default template",
				Subject: "default subject",
			})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the kind has a template", func() {
			BeforeEach(func() {
				template, err := templatesRepo.Create(conn, models.Template{
					ID:      "my-kind-template",
					Name:    "my-kind-template",
					HTML:    "<p>kind template</p>",
					Text:    "some kind template text",
					Subject: "kind subject",
				})
				Expect(err).NotTo(HaveOccurred())

				kind.TemplateID = template.ID
				_, err = kindsRepo.Update(conn, kind)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the template belonging to the kind", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(postal.Templates{
					HTML:    "<p>kind template</p>",
					Text:    "some kind template text",
					Subject: "kind subject",
				}))
			})
		})

		Context("when the client has a template", func() {
			BeforeEach(func() {
				template, err := templatesRepo.Create(conn, models.Template{
					ID:      "my-client-template",
					Name:    "my-client-template",
					HTML:    "<p>client template</p>",
					Text:    "some client template text",
					Subject: "client subject",
				})
				Expect(err).NotTo(HaveOccurred())

				clientsRepo.FindCall.Returns.Client = models.Client{
					ID:         "my-client-id",
					TemplateID: template.ID,
				}
			})

			It("returns the template belonging to the client", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(postal.Templates{
					HTML:    "<p>client template</p>",
					Text:    "some client template text",
					Subject: "client subject",
				}))
			})
		})

		Context("when a templateID is passed", func() {
			BeforeEach(func() {
				templatesCollection.GetCall.Returns.Template = collections.Template{
					Text:     "some testing text",
					Subject:  "some subject",
					HTML:     "<p>v2 awesome</p>",
					ClientID: "my-client-id",
				}
			})

			It("returns the template", func() {
				templates, err := loader.LoadTemplates("my-client-id", "", "some-v2-template-id")
				Expect(err).ToNot(HaveOccurred())

				Expect(templates).To(Equal(postal.Templates{
					HTML:    "<p>v2 awesome</p>",
					Text:    "some testing text",
					Subject: "some subject",
				}))
				Expect(templatesCollection.GetCall.Receives.TemplateID).To(Equal("some-v2-template-id"))
				Expect(templatesCollection.GetCall.Receives.Connection).To(Equal(conn))
				Expect(templatesCollection.GetCall.Receives.ClientID).To(Equal("my-client-id"))
			})
		})

		Context("when the neither client nor kind has a template", func() {
			It("returns the default template", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(postal.Templates{
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
				Expect(templates).To(Equal(postal.Templates{
					HTML:    "<p>The default template</p>",
					Text:    "The default template",
					Subject: "default subject",
				}))
			})
		})

		Context("when the kinds repo has an error", func() {
			BeforeEach(func() {
				kindsRepo.FindError = errors.New("BOOM!")
			})

			It("bubbles up the error", func() {
				_, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).To(HaveOccurred())
			})

		})

		Context("when the clients repo has an error", func() {
			BeforeEach(func() {
				clientsRepo.FindCall.Returns.Error = errors.New("BOOM!")
			})

			It("bubbles up the error", func() {
				_, err := loader.LoadTemplates("my-client-id", "my-kind-id", "")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the templates collection has an error", func() {
			It("returns the error", func() {
				templatesCollection.GetCall.Returns.Error = errors.New("some error on the collection")
				_, err := loader.LoadTemplates("my-client-id", "", "some-v2-template-id")
				Expect(err).To(MatchError("some error on the collection"))
			})
		})
	})
})
