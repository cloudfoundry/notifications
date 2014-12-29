package acceptance

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates Metadata", func() {
	It("Creates a template with metadata", func() {
		// Retrieve Client UAA token
		clientID := "notifications-admin"
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t := TemplateMetadata{
			client:              support.NewClient(Servers.Notifications),
			notificationsServer: Servers.Notifications,
			clientToken:         clientToken,
		}
		t.CreateNewTemplateWithMetadata(support.Template{
			Name:    "Star Wars",
			Subject: "Awesomeness",
			HTML:    "<p>Millenium Falcon</p>",
			Text:    "Millenium Falcon",
			Metadata: map[string]interface{}{
				"some_property": "some_value",
			},
		})
		t.ConfirmMetadataStored(map[string]interface{}{
			"some_property": "some_value",
		})
		t.UpdateTemplateMetadata(support.Template{
			Name:    "Star Wars",
			Subject: "Awesomeness",
			HTML:    "<p>Millenium Falcon</p>",
			Text:    "Millenium Falcon",
			Metadata: map[string]interface{}{
				"hello": true,
			},
		})
		t.ConfirmMetadataStored(map[string]interface{}{
			"hello": true,
		})
	})
})

type TemplateMetadata struct {
	client              *support.Client
	notificationsServer servers.Notifications
	clientToken         uaa.Token
	templateID          string
}

func (t *TemplateMetadata) CreateNewTemplateWithMetadata(template support.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(BeNil())

	t.templateID = templateID
}

func (t *TemplateMetadata) UpdateTemplateMetadata(template support.Template) {
	status, err := t.client.Templates.Update(t.clientToken.Access, t.templateID, template)

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

func (t *TemplateMetadata) ConfirmMetadataStored(metadata map[string]interface{}) {
	status, response, err := t.client.Templates.Get(t.clientToken.Access, t.templateID)

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Metadata).To(Equal(metadata))
}
