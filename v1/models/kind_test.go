package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kind", func() {
	var kind models.Kind

	Describe("TemplateToUse", func() {
		Context("when the template is set", func() {
			BeforeEach(func() {
				kind.TemplateID = "template-id"
			})

			It("returns the template value", func() {
				Expect(kind.TemplateToUse()).To(Equal("template-id"))
			})
		})

		Context("when the template is not set", func() {
			BeforeEach(func() {
				kind.TemplateID = ""
			})

			It("returns the default template value", func() {
				Expect(kind.TemplateToUse()).To(Equal(models.DefaultTemplateID))
			})
		})
	})
})
