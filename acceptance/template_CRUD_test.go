package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create a new template", func() {
	BeforeEach(func() {
		TruncateTables()

		env := config.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		models.NewDatabase(env.DatabaseURL, migrationsPath) // this is the "database" variable
	})

	It("allows a user to create a new template", func() {
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

		test := TemplatesCRUD{}
		name := "Star Wars"
		subject := "Awesomeness"
		html := "<p>Millenium Falcon</p>"
		text := "Millenium Falcon"

		test.CreateNewTemplate(notificationsServer, clientToken, name, subject, html, text)
	})
})

type TemplatesCRUD struct{}

func (test TemplatesCRUD) CreateNewTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, name, subject, html, text string) {
	jsonBody := []byte(fmt.Sprintf(`{"name":"%s", "subject":"%s", "html":"%s", "text":"%s"}`, name, subject, html, text))
	request, err := http.NewRequest("POST", notificationsServer.TemplatesBasePath(), bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusCreated))

	var JSON struct {
		TemplateID string `json:"template-id"`
	}

	err = json.Unmarshal(body, &JSON)
	if err != nil {
		panic(err)
	}

	Expect(JSON.TemplateID).ToNot(BeNil())
}
