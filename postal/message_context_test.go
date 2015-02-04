package postal_test

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageContext", func() {
	var templates postal.Templates
	var email, sender string
	var options postal.Options
	var html postal.HTML
	var delivery postal.Delivery
	var cloak *fakes.Cloak

	BeforeEach(func() {
		email = "bounce@example.com"
		sender = "no-reply@notifications.example.com"

		templates = postal.Templates{
			Text:    "the plainText email < template",
			HTML:    "the html <h1> email < template</h1>",
			Subject: "the subject < template",
		}

		html = postal.HTML{
			BodyContent: "user supplied html",
		}

		options = postal.Options{
			ReplyTo:           "awesomeness",
			Subject:           "the subject",
			KindDescription:   "the kind description",
			SourceDescription: "the source description",
			Text:              "user supplied email text",
			HTML:              html,
			KindID:            "the-kind-id",
			Endorsement:       "this is the endorsement",
			Role:              "OrgRole",
		}

		delivery = postal.Delivery{
			Options:   options,
			UserGUID:  "the-user",
			ClientID:  "the-client-id",
			Email:     email,
			MessageID: "message-id",
			Space: cf.CloudControllerSpace{
				GUID: "my-lovely-guid",
				Name: "the-space",
			},
			Organization: cf.CloudControllerOrganization{
				GUID: "my-super-lovely-guid",
				Name: "the-org",
			},
			Scope: "this.scope",
		}

		cloak = &fakes.Cloak{
			EncryptedResult: []byte("the-encoded-result"),
		}
	})

	Describe("NewMessageContext", func() {
		It("returns the appropriate MessageContext when all options are specified", func() {
			context := postal.NewMessageContext(delivery, sender, cloak, templates)

			Expect(context.From).To(Equal(sender))
			Expect(context.ReplyTo).To(Equal(options.ReplyTo))
			Expect(context.To).To(Equal(email))
			Expect(context.Subject).To(Equal(options.Subject))
			Expect(context.Text).To(Equal(options.Text))
			Expect(context.HTML).To(Equal(options.HTML.BodyContent))
			Expect(context.HTMLComponents).To(Equal(options.HTML))
			Expect(context.TextTemplate).To(Equal(templates.Text))
			Expect(context.HTMLTemplate).To(Equal(templates.HTML))
			Expect(context.SubjectTemplate).To(Equal(templates.Subject))
			Expect(context.KindDescription).To(Equal(options.KindDescription))
			Expect(context.SourceDescription).To(Equal(options.SourceDescription))
			Expect(context.UserGUID).To(Equal("the-user"))
			Expect(context.ClientID).To(Equal("the-client-id"))
			Expect(context.MessageID).To(Equal("message-id"))
			Expect(context.Space).To(Equal("the-space"))
			Expect(context.SpaceGUID).To(Equal("my-lovely-guid"))
			Expect(context.Organization).To(Equal("the-org"))
			Expect(context.OrganizationGUID).To(Equal("my-super-lovely-guid"))
			Expect(context.UnsubscribeID).To(Equal("the-encoded-result"))
			Expect(context.Scope).To(Equal("this.scope"))
			Expect(cloak.DataToEncrypt).To(Equal([]byte("the-user|the-client-id|the-kind-id")))
			Expect(context.Endorsement).To(Equal("this is the endorsement"))
			Expect(context.OrganizationRole).To(Equal("OrgRole"))
		})

		It("falls back to Kind if KindDescription is missing", func() {
			delivery.Options.KindDescription = ""
			context := postal.NewMessageContext(delivery, sender, cloak, templates)

			Expect(context.KindDescription).To(Equal("the-kind-id"))
		})

		It("falls back to clientID when SourceDescription is missing", func() {
			delivery.Options.SourceDescription = ""
			context := postal.NewMessageContext(delivery, sender, cloak, templates)

			Expect(context.SourceDescription).To(Equal("the-client-id"))
		})

		It("fills in subject when subject is not specified", func() {
			delivery.Options.Subject = ""
			context := postal.NewMessageContext(delivery, sender, cloak, templates)
			Expect(context.Subject).To(Equal("[no subject]"))
		})
	})

	Describe("Escape", func() {
		BeforeEach(func() {
			options = postal.Options{
				ReplyTo:           "awesomeness",
				Subject:           "the & subject",
				KindDescription:   "the & kind description",
				SourceDescription: "the & source description",
				Text:              "user & supplied email text",
				HTML:              postal.HTML{BodyContent: "user & supplied html"},
				KindID:            "the & kind",
				Endorsement:       "this & is the endorsement",
				Role:              "OrgRole",
			}

			delivery.Options = options
			delivery.ClientID = "the\"client id"
			delivery.MessageID = "some>id"
			delivery.Space.Name = "the<space"
			delivery.Organization.Name = "the>org"
			delivery.Scope = ""
		})

		It("html escapes various fields on the message context", func() {
			context := postal.NewMessageContext(delivery, sender, cloak, templates)
			context.Escape()

			Expect(context.From).To(Equal("no-reply@notifications.example.com"))
			Expect(context.ReplyTo).To(Equal("awesomeness"))
			Expect(context.To).To(Equal("bounce@example.com"))
			Expect(context.Subject).To(Equal("the &amp; subject"))
			Expect(context.Text).To(Equal("user &amp; supplied email text"))
			Expect(context.HTML).To(Equal("user & supplied html"))
			Expect(context.TextTemplate).To(Equal("the plainText email < template"))
			Expect(context.HTMLTemplate).To(Equal("the html <h1> email < template</h1>"))
			Expect(context.SubjectTemplate).To(Equal("the subject < template"))
			Expect(context.KindDescription).To(Equal("the &amp; kind description"))
			Expect(context.SourceDescription).To(Equal("the &amp; source description"))
			Expect(context.UserGUID).To(Equal("the-user"))
			Expect(context.ClientID).To(Equal("the&#34;client id"))
			Expect(context.MessageID).To(Equal("some&gt;id"))
			Expect(context.Space).To(Equal("the&lt;space"))
			Expect(context.Organization).To(Equal("the&gt;org"))
			Expect(context.Scope).To(Equal(""))
			Expect(context.Endorsement).To(Equal("this &amp; is the endorsement"))
			Expect(context.OrganizationRole).To(Equal("OrgRole"))
		})
	})
})
