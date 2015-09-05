package util_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clock", func() {
	Describe("Now", func() {
		It("should return the current time", func() {
			clock := util.NewClock()

			currentTime := clock.Now()
			Expect(currentTime).To(BeTemporally("~", time.Now()))
		})
	})
})
