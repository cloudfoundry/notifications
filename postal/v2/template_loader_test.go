package v2_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/v2"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("V2TemplateLoader", func() {
	var (
		conn                db.ConnectionInterface
		database            *mocks.Database
		templatesCollection *mocks.TemplatesCollection
		loader              v2.TemplatesLoader
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		templatesCollection = mocks.NewTemplatesCollection()
		loader = v2.NewTemplatesLoader(database, templatesCollection)
	})

	Describe("LoadTemplates", func() {
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

		Context("when the templates collection has an error", func() {
			It("returns the error", func() {
				templatesCollection.GetCall.Returns.Error = errors.New("some error on the collection")

				_, err := loader.LoadTemplates("my-client-id", "", "some-v2-template-id")
				Expect(err).To(MatchError("some error on the collection"))
			})
		})
	})
})
