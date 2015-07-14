package acceptance

import (
	"fmt"
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/warrant"
)

const (
	UAADefaultUsername = "warrant-user"
	UAADefaultClientID = "warrant-client"
)

var (
	UAAHost        string
	UAAAdminClient string
	UAAAdminSecret string
	UAAToken       string
	TraceWriter    io.Writer
)

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
	if os.Getenv("TRACE") == "true" {
		TraceWriter = os.Stdout
	}

	UAAHost = fmt.Sprintf("https://%s", os.Getenv("UAA_HOST"))
	UAAAdminClient = os.Getenv("UAA_ADMIN_CLIENT")
	UAAAdminSecret = os.Getenv("UAA_ADMIN_SECRET")

	warrantClient := warrant.New(warrant.Config{
		Host:          UAAHost,
		SkipVerifySSL: true,
		TraceWriter:   TraceWriter,
	})

	var err error
	UAAToken, err = warrantClient.Clients.GetToken(UAAAdminClient, UAAAdminSecret)
	Expect(err).NotTo(HaveOccurred())
})
