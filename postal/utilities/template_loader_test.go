package utilities_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateLoader", func() {
	var loader utilities.TemplatesLoader
	var finder *fakes.TemplateFinder
	var clientsRepo *fakes.ClientsRepo
	var kindsRepo *fakes.KindsRepo
	var templatesRepo *fakes.TemplatesRepo
	var conn models.ConnectionInterface
	var database *fakes.Database

	BeforeEach(func() {
		finder = fakes.NewTemplateFinder()
		clientsRepo = fakes.NewClientsRepo()
		kindsRepo = fakes.NewKindsRepo()
		templatesRepo = fakes.NewTemplatesRepo()
		database = fakes.NewDatabase()
		conn = database.Connection()
		loader = utilities.NewTemplatesLoader(finder, database, clientsRepo, kindsRepo, templatesRepo)
	})

	Describe("LoadTemplates", func() {
		var kind models.Kind
		var client models.Client

		BeforeEach(func() {
			var err error

			client, err = clientsRepo.Create(conn, models.Client{
				ID: "my-client-id",
			})
			if err != nil {
				panic(err)
			}

			kind, err = kindsRepo.Create(conn, models.Kind{
				ID:       "my-kind-id",
				ClientID: client.ID,
			})
			if err != nil {
				panic(err)
			}

			_, err = templatesRepo.Create(conn, models.Template{
				ID:      "default",
				Name:    "Default Template",
				HTML:    "<p>The default template</p>",
				Text:    "The default template",
				Subject: "default subject",
			})
			if err != nil {
				panic(err)
			}
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
				if err != nil {
					panic(err)
				}

				kind.Template = template.ID
				_, err = kindsRepo.Update(conn, kind)
				if err != nil {
					panic(err)
				}
			})

			It("returns the template belonging to the kind", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", models.UserBodyTemplateName, models.SubjectProvidedTemplateName)
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
				if err != nil {
					panic(err)
				}

				client.Template = template.ID
				_, err = clientsRepo.Update(conn, client)
				if err != nil {
					panic(err)
				}
			})

			It("returns the template belonging to the client", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", models.UserBodyTemplateName, models.SubjectProvidedTemplateName)
				Expect(err).ToNot(HaveOccurred())
				Expect(templates).To(Equal(postal.Templates{
					HTML:    "<p>client template</p>",
					Text:    "some client template text",
					Subject: "client subject",
				}))
			})
		})

		Context("when the neither client nor kind has a template", func() {
			It("returns the default template", func() {
				templates, err := loader.LoadTemplates("my-client-id", "my-kind-id", models.UserBodyTemplateName, models.SubjectProvidedTemplateName)
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
				_, err := loader.LoadTemplates("my-client-id", "my-kind-id", models.UserBodyTemplateName, models.SubjectProvidedTemplateName)
				Expect(err).To(HaveOccurred())
			})

		})

		Context("when the clients repo has an error", func() {
			BeforeEach(func() {
				clientsRepo.FindError = errors.New("BOOM!")
			})

			It("bubbles up the error", func() {
				_, err := loader.LoadTemplates("my-client-id", "my-kind-id", models.UserBodyTemplateName, models.SubjectProvidedTemplateName)
				Expect(err).To(HaveOccurred())
			})
		})

	})
})
