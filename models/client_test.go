package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var client models.Client

	Describe("TemplateToUse", func() {
		Context("when the template is set", func() {
			BeforeEach(func() {
				client.Template = "template-id"
			})

			It("returns the template value", func() {
				Expect(client.TemplateToUse()).To(Equal("template-id"))
			})
		})

		Context("when the template is not set", func() {
			BeforeEach(func() {
				client.Template = ""
			})

			It("returns the default template value", func() {
				Expect(client.TemplateToUse()).To(Equal(models.DefaultTemplateID))
			})
		})
	})
})
