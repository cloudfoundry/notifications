package docs_test

import (
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/notifications/docs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestInspector", func() {
	Describe("GetResourceInfo", func() {
		Context("when the resource is a top-level list", func() {
			var resourceInfo docs.ResourceInfo

			BeforeEach(func() {
				url, err := url.Parse("/senders")
				Expect(err).NotTo(HaveOccurred())
				request := &http.Request{
					URL: url,
				}

				inspector := docs.NewRequestInspector()
				resourceInfo = inspector.GetResourceInfo(request)
			})

			It("should return the cannonical resource type from the request's URL", func() {

				Expect(resourceInfo.ResourceType).To(Equal("senders"))
				Expect(resourceInfo.ListName).To(Equal("Senders"))
				Expect(resourceInfo.ItemName).To(Equal("Sender"))
			})

			It("sets a flag for the list resource", func() {
				Expect(resourceInfo.IsItem).To(BeFalse())
			})
		})

		Context("when the resource is a top-level item", func() {
			var resourceInfo docs.ResourceInfo

			BeforeEach(func() {
				url, err := url.Parse("/senders/2daa96a8-5e58-11e5-8c91-1b8480f7bd21")
				Expect(err).NotTo(HaveOccurred())
				request := &http.Request{
					URL: url,
				}

				inspector := docs.NewRequestInspector()
				resourceInfo = inspector.GetResourceInfo(request)
			})

			It("should return the cannonical resource type from the request's URL", func() {
				Expect(resourceInfo.ResourceType).To(Equal("senders"))
				Expect(resourceInfo.ListName).To(Equal("Senders"))
				Expect(resourceInfo.ItemName).To(Equal("Sender"))
			})

			It("sets a flag for the item resource", func() {
				Expect(resourceInfo.IsItem).To(BeTrue())
			})
		})

		Context("when the resource is a nested list", func() {
			It("should return the cannonical resource type from the request's URL", func() {
				inspector := docs.NewRequestInspector()

				url, err := url.Parse("/senders/2daa96a8-5e58-11e5-8c91-1b8480f7bd21/campaign_types")
				Expect(err).NotTo(HaveOccurred())
				request := &http.Request{
					URL: url,
				}
				resourceInfo := inspector.GetResourceInfo(request)

				Expect(resourceInfo.ResourceType).To(Equal("campaign_types"))
				Expect(resourceInfo.ListName).To(Equal("Campaign types"))
				Expect(resourceInfo.ItemName).To(Equal("Campaign type"))
			})
		})

		Context("when the URL contains a hostname", func() {
			It("should return the cannonical resource type from the request's URL", func() {
				inspector := docs.NewRequestInspector()

				url, err := url.Parse("https://some.example.com:123/senders/2daa96a8-5e58-11e5-8c91-1b8480f7bd21")
				Expect(err).NotTo(HaveOccurred())
				request := &http.Request{
					URL: url,
				}
				resourceInfo := inspector.GetResourceInfo(request)

				Expect(resourceInfo.ResourceType).To(Equal("senders"))
				Expect(resourceInfo.ListName).To(Equal("Senders"))
				Expect(resourceInfo.ItemName).To(Equal("Sender"))
			})
		})
	})
})
