package network_test

import (
	"io/ioutil"
	"net/url"

	"github.com/pivotal-cf-experimental/warrant/internal/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestBodyEncoder", func() {
	Describe("JSONRequestBody", func() {
		Describe("Encode", func() {
			It("returns a JSON encoded representation of the given object with proper content type", func() {
				var object struct {
					Hello string `json:"hello"`
				}
				object.Hello = "goodbye"

				body, contentType, err := network.NewJSONRequestBody(object).Encode()
				Expect(err).NotTo(HaveOccurred())
				Expect(ioutil.ReadAll(body)).To(MatchJSON(`{
					"hello": "goodbye"
				}`))
				Expect(contentType).To(Equal("application/json"))
			})

			It("returns an error when the JSON cannot be encoded", func() {
				_, _, err := network.NewJSONRequestBody(func() {}).Encode()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FormRequestBody", func() {
		Describe("Encode", func() {
			It("returns a form URL encoded representation of the given object with proper content type", func() {
				values := url.Values{
					"hello": []string{"goodbye"},
					"black": []string{"white"},
				}

				body, contentType, err := network.NewFormRequestBody(values).Encode()
				Expect(err).NotTo(HaveOccurred())
				Expect(ioutil.ReadAll(body)).To(BeEquivalentTo("black=white&hello=goodbye"))
				Expect(contentType).To(Equal("application/x-www-form-urlencoded"))
			})
		})
	})
})
