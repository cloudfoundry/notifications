package acceptance

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/nu7hatch/gouuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var index int

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func NewGUID(prefix string) string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	if prefix != "" {
		return fmt.Sprintf("%s-%s", prefix, guid.String())
	}

	return guid.String()
}

func NewOrgName(prefix string) string {
	index++

	if prefix != "" {
		return fmt.Sprintf("%s-%d", prefix, index)
	}

	return strconv.Itoa(index)
}
