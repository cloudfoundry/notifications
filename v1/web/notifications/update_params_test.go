package notifications_test

import (
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification", func() {
	Describe("NewNotificationParams", func() {
		Context("when the json is valid", func() {
			It("returns a NotificationUpdateParams", func() {
				body := strings.NewReader(`{"description":"my awesome notification", "critical":true, "template":"my-awesome-template"}`)
				updateParams, err := notifications.NewNotificationParams(body)
				Expect(err).ToNot(HaveOccurred())
				Expect(updateParams).To(BeAssignableToTypeOf(notifications.NotificationUpdateParams{}))
			})
		})

		Context("when the json is invalid", func() {
			Context("when the json is missing a required field", func() {
				It("returns a validation error", func() {
					body := strings.NewReader(`{"critical":true, "template":"my-awesome-template"}`)
					_, err := notifications.NewNotificationParams(body)
					Expect(err).To(BeAssignableToTypeOf(webutil.ValidationError{}))
				})
			})

			Context("when the json is malformed", func() {
				It("returns a parse error", func() {
					body := strings.NewReader(`{"description":"my awesome notification", "critical":true, "template":"my-awesome-template}`)
					_, err := notifications.NewNotificationParams(body)
					Expect(err).To(BeAssignableToTypeOf(webutil.ParseError{}))
				})
			})
		})
	})

	Describe("ToModel", func() {
		It("returns a model.Kind composed of the NotificationUpdateParams", func() {
			body := strings.NewReader(`{"description":"my awesome notification", "critical":true, "template":"my-awesome-template"}`)
			updateParams, err := notifications.NewNotificationParams(body)
			Expect(err).NotTo(HaveOccurred())

			notification := updateParams.ToModel("client-id", "notification-id")
			Expect(notification.Description).To(Equal("my awesome notification"))
			Expect(notification.Critical).To(Equal(true))
			Expect(notification.TemplateID).To(Equal("my-awesome-template"))
			Expect(notification.ClientID).To(Equal("client-id"))
			Expect(notification.ID).To(Equal("notification-id"))
		})
	})
})
