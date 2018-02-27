package clients_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

func TestClientsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/server/clients")
}

func NewTokens() *domain.Tokens {
	return domain.NewTokens(common.TestPublicKey, common.TestPrivateKey, []string{})
}
