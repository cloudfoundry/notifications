package common_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageContext", func() {
	var templates common.Templates
	var email, sender, domain string
	var options common.Options
	var html common.HTML
	var delivery common.Delivery
	var cloak *mocks.Cloak
	var reqReceived time.Time

	BeforeEach(func() {
		email = "bounce@example.com"
		sender = "no-reply@notifications.example.com"
		domain = "http://www.example.com"

		templates = common.Templates{
			Text:    "the plainText email < template",
			HTML:    "the html <h1> email < template</h1>",
			Subject: "the subject < template",
		}

		html = common.HTML{
			BodyContent: "user supplied html",
		}

		options = common.Options{
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

		reqReceived, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:40:12.207187819-07:00")

		delivery = common.Delivery{
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
			Scope:           "this.scope",
			RequestReceived: reqReceived,
		}

		cloak = mocks.NewCloak()
		cloak.VeilCall.Returns.CipherText = []byte("the-encoded-result")
	})

	Describe("NewMessageContext", func() {
		It("returns the appropriate MessageContext when all options are specified", func() {
			context := common.NewMessageContext(delivery, sender, domain, cloak, templates)

			Expect(cloak.VeilCall.Receives.PlainText).To(Equal([]byte("the-user|the-client-id|the-kind-id")))

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
			Expect(context.Endorsement).To(Equal("this is the endorsement"))
			Expect(context.OrganizationRole).To(Equal("OrgRole"))
			Expect(context.RequestReceived).To(Equal(reqReceived))
			Expect(context.Domain).To(Equal(domain))
		})

		It("falls back to Kind if KindDescription is missing", func() {
			delivery.Options.KindDescription = ""
			context := common.NewMessageContext(delivery, sender, domain, cloak, templates)

			Expect(context.KindDescription).To(Equal("the-kind-id"))
		})

		It("falls back to clientID when SourceDescription is missing", func() {
			delivery.Options.SourceDescription = ""
			context := common.NewMessageContext(delivery, sender, domain, cloak, templates)

			Expect(context.SourceDescription).To(Equal("the-client-id"))
		})

		It("fills in subject when subject is not specified", func() {
			delivery.Options.Subject = ""
			context := common.NewMessageContext(delivery, sender, domain, cloak, templates)
			Expect(context.Subject).To(Equal("[no subject]"))
		})
	})

	Describe("Escape", func() {
		BeforeEach(func() {
			options = common.Options{
				ReplyTo:           "awesomeness",
				Subject:           "the & subject",
				KindDescription:   "the & kind description",
				SourceDescription: "the & source description",
				Text:              "user & supplied email text",
				HTML:              common.HTML{BodyContent: "user & supplied html"},
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
			context := common.NewMessageContext(delivery, sender, domain, cloak, templates)
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
