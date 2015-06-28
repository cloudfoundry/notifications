package params_test

import (
	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	Describe("EmailValidator", func() {
		var (
			notify    *params.NotifyParams
			validator params.EmailValidator
		)

		BeforeEach(func() {
			notify = &params.NotifyParams{
				Text: "my silly text",
				To:   "bob@example.com",
			}
			validator = params.EmailValidator{}
		})

		Describe("Validate", func() {
			It("validates the email fields on Notify", func() {
				Expect(validator.Validate(notify)).To(BeTrue())
				Expect(len(notify.Errors)).To(Equal(0))

				notify.To = ""

				Expect(validator.Validate(notify)).To(BeFalse())
				Expect(len(notify.Errors)).To(Equal(1))
				Expect(notify.Errors).To(ContainElement(`"to" is a required field`))

				notify.Text = ""

				Expect(validator.Validate(notify)).To(BeFalse())
				Expect(len(notify.Errors)).To(Equal(2))
				Expect(notify.Errors).To(ContainElement(`"to" is a required field`))
				Expect(notify.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

				notify.To = "otherUser@example.com"
				notify.ParsedHTML = params.HTML{BodyContent: "<p>Contents of this email message</p>"}

				Expect(validator.Validate(notify)).To(BeTrue())
				Expect(len(notify.Errors)).To(Equal(0))
			})

			Context("When the notify params object finds an invalid email", func() {
				It("Reports a validation error", func() {
					notify.To = params.InvalidEmail

					Expect(validator.Validate(notify)).To(BeFalse())
					Expect(len(notify.Errors)).To(Equal(1))
					Expect(notify.Errors).To(ContainElement(`"to" is improperly formatted`))
				})
			})
		})
	})

	Describe("GUIDValidator", func() {
		var (
			notify    *params.NotifyParams
			validator params.GUIDValidator
		)

		BeforeEach(func() {
			notify = &params.NotifyParams{
				KindID:  "test_email",
				Subject: "Summary of contents",
				Text:    "Contents of the email message",
			}
			validator = params.GUIDValidator{}
		})

		Describe("Validate", func() {
			It("validates the kind and text fields", func() {
				Expect(validator.Validate(notify)).To(BeTrue())
				Expect(len(notify.Errors)).To(Equal(0))

				notify.KindID = ""

				Expect(validator.Validate(notify)).To(BeFalse())
				Expect(len(notify.Errors)).To(Equal(1))
				Expect(notify.Errors).To(ContainElement(`"kind_id" is a required field`))

				notify.Text = ""

				Expect(validator.Validate(notify)).To(BeFalse())
				Expect(len(notify.Errors)).To(Equal(2))
				Expect(notify.Errors).To(ContainElement(`"kind_id" is a required field`))
				Expect(notify.Errors).To(ContainElement(`"text" or "html" fields must be supplied`))

				notify.KindID = "something"
				notify.ParsedHTML.BodyContent = "<p>banana</p>"

				Expect(validator.Validate(notify)).To(BeTrue())
				Expect(len(notify.Errors)).To(Equal(0))
			})

			It("validates that KindID is properly formatted", func() {
				notify.KindID = "A_valid.id-99"

				Expect(validator.Validate(notify)).To(BeTrue())
				Expect(len(notify.Errors)).To(Equal(0))

				notify.KindID = "an_invalid.id-00!"

				Expect(validator.Validate(notify)).To(BeFalse())
				Expect(len(notify.Errors)).To(Equal(1))
				Expect(notify.Errors).To(ContainElement(`"kind_id" is improperly formatted`))
			})

			It("validates that the role must be OrgManager, OrgAuditor, BillingManager, or empty", func() {
				for _, role := range []string{"OrgManager", "OrgAuditor", "BillingManager", ""} {
					notify.Role = role
					Expect(validator.Validate(notify)).To(BeTrue())
					Expect(len(notify.Errors)).To(Equal(0))
				}

				notify.Role = "bad-role-name"
				Expect(validator.Validate(notify)).To(BeFalse())
				Expect(len(notify.Errors)).To(Equal(1))
				Expect(notify.Errors).To(ContainElement(`"role" must be "OrgManager", "OrgAuditor", "BillingManager" or unset`))
			})
		})
	})
})
