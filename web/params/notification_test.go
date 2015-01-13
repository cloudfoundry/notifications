package params_test

import (
	"strings"

	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification", func() {
	Describe("NewNotificationParams", func() {
		Context("when the json is valid", func() {
			It("returns a NotificationUpdateParams", func() {
				body := strings.NewReader(`{"description":"my awesome notification", "critical":true, "template":"my-awesome-template"}`)
				updateParams, err := params.NewNotificationParams(body)
				Expect(err).ToNot(HaveOccurred())
				Expect(updateParams).To(BeAssignableToTypeOf(params.NotificationUpdateParams{}))
			})
		})

		Context("when the json is invalid", func() {
			Context("when the json is missing a required field", func() {
				It("returns a validation error", func() {
					body := strings.NewReader(`{"critical":true, "template":"my-awesome-template"}`)
					_, err := params.NewNotificationParams(body)
					Expect(err).To(BeAssignableToTypeOf(params.ValidationError{}))
				})
			})

			Context("when the json is malformed", func() {
				It("returns a parse error", func() {
					body := strings.NewReader(`{"description":"my awesome notification", "critical":true, "template":"my-awesome-template}`)
					_, err := params.NewNotificationParams(body)
					Expect(err).To(BeAssignableToTypeOf(params.ParseError{}))
				})
			})
		})
	})

	Describe("ToModel", func() {
		It("returns a model.Kind composed of the NotificationUpdateParams", func() {
			body := strings.NewReader(`{"description":"my awesome notification", "critical":true, "template":"my-awesome-template"}`)
			updateParams, err := params.NewNotificationParams(body)
			if err != nil {
				panic(err)
			}

			notification := updateParams.ToModel()
			Expect(notification.Description).To(Equal("my awesome notification"))
			Expect(notification.Critical).To(Equal(true))
			Expect(notification.TemplateID).To(Equal("my-awesome-template"))
		})
	})
})
