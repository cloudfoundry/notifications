package util_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type errorReader struct{}

func (e errorReader) Read([]byte) (int, error) {
	return 0, errors.New("failed to read")
}

var _ = Describe("IDGenerator", func() {
	It("generates a ID without generating an error", func() {
		reader := bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz1234567890"))
		generator := util.NewIDGenerator(reader)

		guid, err := generator.Generate()
		Expect(err).NotTo(HaveOccurred())
		Expect(guid).To(Equal("61626364-6566-6768-696a-6b6c6d6e6f70"))
	})

	It("returns an error if the reader errors", func() {
		reader := errorReader{}
		generator := util.NewIDGenerator(reader)

		_, err := generator.Generate()
		Expect(err).To(MatchError(errors.New("failed to read")))
	})
})
