package models_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type errorReader struct{}

func (e errorReader) Read([]byte) (int, error) {
	return 0, errors.New("failed to read")
}

var _ = Describe("GUIDGenerator", func() {
	It("generates a GUID without generating an error", func() {
		reader := bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz1234567890"))
		generator := models.NewGUIDGenerator(reader)

		guid, err := generator.Generate()
		Expect(err).NotTo(HaveOccurred())
		Expect(guid).To(Equal("61626364-6566-6768-696a-6b6c6d6e6f70"))
	})

	It("returns an error if the reader errors", func() {
		reader := errorReader{}
		generator := models.NewGUIDGenerator(reader)

		_, err := generator.Generate()
		Expect(err).To(MatchError(errors.New("failed to read")))
	})
})
