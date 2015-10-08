package acceptance

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/docs"
	"github.com/cloudfoundry-incubator/notifications/testing/servers"
	"github.com/cloudfoundry-incubator/notifications/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/pivotal-cf-experimental/warrant"
	"github.com/pivotal-cf-experimental/warrant/testserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const PUBLIC_KEY = `banana`

var (
	Servers struct {
		Notifications servers.Notifications
		SMTP          *servers.SMTP
		UAA           *testserver.UAA
		CC            servers.CC
	}
	Trace, _ = strconv.ParseBool(os.Getenv("TRACE"))

	adminToken        string
	warrantClient     warrant.Warrant
	roundtripRecorder *docs.RoundTripRecorder
	users             map[string]string
)

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/acceptance")
}

var _ = BeforeSuite(func() {
	Servers.SMTP = servers.NewSMTP()
	Servers.SMTP.Boot()

	Servers.UAA = testserver.NewUAA(testserver.Config{
		PublicKey: PUBLIC_KEY,
	})
	Servers.UAA.Start()
	os.Setenv("UAA_HOST", Servers.UAA.URL())

	// Create the notfications client
	adminToken = Servers.UAA.ClientTokenFor("admin", []string{"clients.write", "clients.read", "scim.write", "password.write"}, []string{"clients", "scim", "password"})
	var traceWriter io.Writer
	if os.Getenv("TRACE") == "true" {
		traceWriter = os.Stdout
	}
	warrantClient = warrant.New(warrant.Config{
		Host:        Servers.UAA.URL(),
		TraceWriter: traceWriter,
	})

	err := warrantClient.Clients.Create(warrant.Client{
		ID:    os.Getenv("UAA_CLIENT_ID"),
		Scope: []string{"cloud_controller.admin", "scim.read"},
	}, os.Getenv("UAA_CLIENT_SECRET"), adminToken)
	Expect(err).NotTo(HaveOccurred())

	testUsers := map[string]string{
		"user-123":          "user-123@example.com",
		"user-456":          "user-456@example.com",
		"user-789":          "user-789",
		"unauthorized-user": "unauthorized-user@example.com",
	}

	users = map[string]string{}
	for userName, email := range testUsers {
		user, err := warrantClient.Users.Create(userName, email, adminToken)
		Expect(err).NotTo(HaveOccurred())

		err = warrantClient.Users.SetPassword(user.ID, "password", adminToken)
		Expect(err).NotTo(HaveOccurred())

		users[userName] = user.ID
	}

	Servers.CC = servers.NewCC(users)
	Servers.CC.Boot()

	Servers.Notifications = servers.NewNotifications()
	Servers.Notifications.Compile()
	Servers.Notifications.ResetDatabase()
	Servers.Notifications.Boot()

	roundtripRecorder = docs.NewRoundTripRecorder()
})

var _ = AfterSuite(func() {
	Servers.Notifications.Close()
	Servers.Notifications.Destroy()

	Servers.UAA.Close()

	if roundtripRecorder != nil {
		context, err := docs.BuildTemplateContext(docs.Structure, roundtripRecorder.RoundTrips)
		Expect(err).NotTo(HaveOccurred())

		markdown, err := docs.GenerateMarkdown(context)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(fmt.Sprintf("%s/V2_API.md", os.Getenv("ROOT_PATH")), []byte(markdown), 0644)
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = BeforeEach(func() {
	Servers.Notifications.ResetDatabase()
	Servers.SMTP.Reset()
})

var _ = AfterEach(func() {
	err := Servers.Notifications.WaitForJobsQueueToEmpty()
	Expect(err).NotTo(HaveOccurred())
})

func GetClientTokenWithScopes(scopes ...string) (string, error) {
	id, err := util.NewIDGenerator(rand.Reader).Generate()
	if err != nil {
		return "", err
	}

	err = warrantClient.Clients.Create(warrant.Client{
		ID:    id,
		Scope: scopes,
	}, "secret", adminToken)
	if err != nil {
		return "", err
	}

	token, err := warrantClient.Clients.GetToken(id, "secret")
	if err != nil {
		return "", err
	}

	_, err = jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(PUBLIC_KEY), nil
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func UpdateClientTokenWithDifferentScopes(token string, scopes ...string) (string, error) {
	t, err := warrantClient.Tokens.Decode(token)
	if err != nil {
		return "", err
	}

	client, err := warrantClient.Clients.Get(t.ClientID, adminToken)
	if err != nil {
		return "", err
	}

	client.Scope = scopes
	err = warrantClient.Clients.Update(client, adminToken)
	if err != nil {
		return "", err
	}

	token, err = warrantClient.Clients.GetToken(client.ID, "secret")
	if err != nil {
		return "", err
	}

	_, err = jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(PUBLIC_KEY), nil
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserTokenAndIdFor(userName string) (string, string, error) {
	userGUID := users[userName]
	token, err := warrantClient.Users.GetToken(userName, "password")
	if err != nil {
		return "", "", err
	}

	_, err = jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(PUBLIC_KEY), nil
	})
	if err != nil {
		return "", "", err
	}

	return token, userGUID, nil
}
