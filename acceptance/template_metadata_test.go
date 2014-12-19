package acceptance

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates Metadata", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("Creates a template with metadata", func() {
		// Boot Fake SMTP Server
		smtpServer := servers.NewSMTP()
		smtpServer.Boot()

		// Boot Fake UAA Server
		uaaServer := servers.NewUAA()
		uaaServer.Boot()
		defer uaaServer.Close()

		// Boot Real Notifications Server
		notificationsServer := servers.NewNotifications()
		notificationsServer.Boot()
		defer notificationsServer.Close()

		// Retrieve Client UAA token
		clientID := "notifications-admin"
		env := config.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		test := TemplateMetadata{
			client:              support.NewClient(notificationsServer),
			notificationsServer: notificationsServer,
			clientToken:         clientToken,
		}

		createdTemplate := params.Template{
			Name:    "Star Wars",
			Subject: "Awesomeness",
			HTML:    "<p>Millenium Falcon</p>",
			Text:    "Millenium Falcon",
			Metadata: map[string]interface{}{
				"some_property": "some_value",
			},
		}

		test.CreateNewTemplateWithMetadata(createdTemplate)
		test.GetTemplateWithMetadata()
	})
})

type TemplateMetadata struct {
	client              *support.Client
	notificationsServer servers.Notifications
	clientToken         uaa.Token
	templateID          string
}

func (test *TemplateMetadata) CreateNewTemplateWithMetadata(template params.Template) {
	status, templateID, err := test.client.Templates.Create(test.clientToken.Access, template)

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(BeNil())

	test.templateID = templateID
}

func (test *TemplateMetadata) GetTemplateWithMetadata() {
	status, response, err := test.client.Templates.Get(test.clientToken.Access, test.templateID)

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Metadata).To(Equal(map[string]interface{}{
		"some_property": "some_value",
	}))
}
