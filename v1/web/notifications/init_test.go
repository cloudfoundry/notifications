package notifications_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWebV1NotificationsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/web/notifications")
}
