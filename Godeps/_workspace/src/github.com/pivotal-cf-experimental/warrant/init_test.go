package warrant_test

import (
	"io"
	"os"
	"testing"

	"github.com/pivotal-cf-experimental/warrant/internal/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	fakeUAA          *server.UAA
	fakeUAAPublicKey string
	TraceWriter      io.Writer
)

func TestWarrantSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Warrant Suite")
}

var _ = BeforeSuite(func() {
	if os.Getenv("TRACE") == "true" {
		TraceWriter = os.Stdout
	}

	fakeUAAPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0m59l2u9iDnMbrXHfqkO
rn2dVQ3vfBJqcDuFUK03d+1PZGbVlNCqnkpIJ8syFppW8ljnWweP7+LiWpRoz0I7
fYb3d8TjhV86Y997Fl4DBrxgM6KTJOuE/uxnoDhZQ14LgOU2ckXjOzOdTsnGMKQB
LCl0vpcXBtFLMaSbpv1ozi8h7DJyVZ6EnFQZUWGdgTMhDrmqevfx95U/16c5WBDO
kqwIn7Glry9n9Suxygbf8g5AzpWcusZgDLIIZ7JTUldBb8qU2a0Dl4mvLZOn4wPo
jfj9Cw2QICsc5+Pwf21fP+hzf+1WSRHbnYv8uanRO0gZ8ekGaghM/2H6gqJbo2nI
JwIDAQAB
-----END PUBLIC KEY-----`

	fakeUAA = server.NewUAA(server.Config{
		PublicKey: fakeUAAPublicKey,
	})
	fakeUAA.Start()
})

var _ = AfterSuite(func() {
	fakeUAA.Close()
})

var _ = BeforeEach(func() {
	fakeUAA.Reset()
})
