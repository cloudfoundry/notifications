package documents_test

import (
	"encoding/json"
	"time"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Meta", func() {
	It("marshals itself into JSON with the correct timestamp format", func() {
		created, err := time.Parse(time.RFC1123, "Thu, 19 Mar 2015 11:59:05 UTC")
		Expect(err).NotTo(HaveOccurred())

		lastModified := created.Add(25 * time.Second)

		output, err := json.Marshal(documents.Meta{
			Version:      10,
			Created:      created,
			LastModified: lastModified,
		})

		Expect(output).To(MatchJSON(`{
			"version": 10,
			"created": "2015-03-19T11:59:05.000Z",
			"lastModified": "2015-03-19T11:59:30.000Z"
		}`))
	})
})
