package notify_test

import (
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	Describe("EmailValidator", func() {
		var (
			params    *notify.NotifyParams
			validator notify.EmailValidator
		)

		BeforeEach(func() {
			params = &notify.NotifyParams{
				Text: "my silly text",
				To:   "bob@example.com",
			}
			validator = notify.EmailValidator{}
		})

		Describe("Validate", func() {
			It("validates the email fields on Notify", func() {
				Expect(validator.Validate(params)).To(BeTrue())
				Expect(len(params.Errors)).To(Equal(0))

				params.To = ""

				Expect(validator.Validate(params)).To(BeFalse())
				Expect(len(params.Errors)).To(Equal(1))
				Expect(params.Errors).To(ContainElement(`"to" is a required field`))

				params.Text = ""

				Expect(validator.Validate(params)).To(BeFalse())
				Expect(len(params.Errors)).To(Equal(2))
				Expect(params.Errors).To(ContainElement(`"to" is a required field`))
				Expect(params.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

				params.To = "otherUser@example.com"
				params.ParsedHTML = notify.HTML{BodyContent: "<p>Contents of this email message</p>"}

				Expect(validator.Validate(params)).To(BeTrue())
				Expect(len(params.Errors)).To(Equal(0))
			})

			Context("When the notify params object finds an invalid email", func() {
				It("Reports a validation error", func() {
					params.To = notify.InvalidEmail

					Expect(validator.Validate(params)).To(BeFalse())
					Expect(len(params.Errors)).To(Equal(1))
					Expect(params.Errors).To(ContainElement(`"to" is improperly formatted`))
				})
			})
		})
	})

	Describe("GUIDValidator", func() {
		var (
			params    *notify.NotifyParams
			validator notify.GUIDValidator
		)

		BeforeEach(func() {
			params = &notify.NotifyParams{
				KindID:  "test_email",
				Subject: "Summary of contents",
				Text:    "Contents of the email message",
			}
			validator = notify.GUIDValidator{}
		})

		Describe("Validate", func() {
			It("validates the kind and text fields", func() {
				Expect(validator.Validate(params)).To(BeTrue())
				Expect(len(params.Errors)).To(Equal(0))

				params.KindID = ""

				Expect(validator.Validate(params)).To(BeFalse())
				Expect(len(params.Errors)).To(Equal(1))
				Expect(params.Errors).To(ContainElement(`"kind_id" is a required field`))

				params.Text = ""

				Expect(validator.Validate(params)).To(BeFalse())
				Expect(len(params.Errors)).To(Equal(2))
				Expect(params.Errors).To(ContainElement(`"kind_id" is a required field`))
				Expect(params.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

				params.KindID = "something"
				params.ParsedHTML.BodyContent = "<p>banana</p>"

				Expect(validator.Validate(params)).To(BeTrue())
				Expect(len(params.Errors)).To(Equal(0))
			})

			It("validates that KindID is properly formatted", func() {
				params.KindID = "A_valid.id-99"

				Expect(validator.Validate(params)).To(BeTrue())
				Expect(len(params.Errors)).To(Equal(0))

				params.KindID = "an_invalid.id-00!"

				Expect(validator.Validate(params)).To(BeFalse())
				Expect(len(params.Errors)).To(Equal(1))
				Expect(params.Errors).To(ContainElement(`"kind_id" is improperly formatted`))
			})

			It("validates that the role must be OrgManager, OrgAuditor, BillingManager, or empty", func() {
				for _, role := range []string{"OrgManager", "OrgAuditor", "BillingManager", ""} {
					params.Role = role
					Expect(validator.Validate(params)).To(BeTrue())
					Expect(len(params.Errors)).To(Equal(0))
				}

				params.Role = "bad-role-name"
				Expect(validator.Validate(params)).To(BeFalse())
				Expect(len(params.Errors)).To(Equal(1))
				Expect(params.Errors).To(ContainElement(`"role" must be "OrgManager", "OrgAuditor", "BillingManager" or unset`))
			})
		})
	})
})
