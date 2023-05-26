package docs_test

import (
	"github.com/cloudfoundry-incubator/notifications/docs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diff", func() {
	Context("when there is a unimportant difference in the Date header", func() {
		It("returns false", func() {
			left := "Date: Fri, 09 Oct 2015 16:26:03 GMT"
			right := "Date: Thu, 08 Oct 2015 12:12:03 PST"

			Expect(docs.Diff(left, right)).To(BeFalse())
		})
	})

	Context("when there is a real difference in the Date header", func() {
		It("returns false", func() {
			left := "Date: Fri, 09 Oct 2015 16:26:03 GMT"
			right := "Date: Thu, 08 Oct 2015 12:12:03 PST banana"

			Expect(docs.Diff(left, right)).To(BeTrue())
		})
	})

	Context("when there is an unimportant difference in the Authorization header", func() {
		It("returns false", func() {
			left := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQi\nsomething"
			right := "Authorization: Bearer EYjHBgCIoIjiuZi1nIiSiNr5Cci6iKPxvcj9.EYjHDwqI\nsomething"

			Expect(docs.Diff(left, right)).To(BeFalse())
		})
	})

	Context("when there is an important difference in the Authorization header", func() {
		It("returns true", func() {
			left := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQi"
			right := "Authorization: Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQi"

			Expect(docs.Diff(left, right)).To(BeTrue())
		})

		It("returns true", func() {
			left := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQi\nsomething"
			right := "Authorization: Bearer EYjHBgCIoIjiuZi1nIiSiNr5Cci6iKPxvcj9.EYjHDwqI\nsomething hello"

			Expect(docs.Diff(left, right)).To(BeTrue())
		})
	})

	Context("when there is an unimportant change to a guid", func() {
		It("returns false", func() {
			left := "/campaigns/8bba8e63-41e8-3cca-c098-da8c3500deac/status"
			right := "/campaigns/8e69a00a-da26-36ed-b4fc-28a2150c6593/status"

			Expect(docs.Diff(left, right)).To(BeFalse())
		})

		It("returns false", func() {
			left := "\"id\": \"8bba8e63-41e8-3cca-c098-da8c3500deac\""
			right := "\"id\": \"dda6a0c0-669b-1572-7267-67971a126dad\""

			Expect(docs.Diff(left, right)).To(BeFalse())
		})
	})

	Context("when there is an important change to a guid-relate string", func() {
		It("returns true", func() {
			left := "/campaigns/8bba8e63-41e8-3cca-c098-da8c3500deac/status"
			right := "/senders/8e69a00a-da26-36ed-b4fc-28a2150c6593"

			Expect(docs.Diff(left, right)).To(BeTrue())
		})

		It("returns true", func() {
			left := "/campaigns/8bba8e63-41e8-3cca-c098-da8c3500deac/status"
			right := "/campaigns/8bbae63-41e8-3cca-c098-da8c3500deac/status"

			Expect(docs.Diff(left, right)).To(BeTrue())
		})
	})

	Context("when there is an unimportant change to a timestamp", func() {
		It("returns false", func() {
			left := "2015-10-13T22:19:44Z"
			right := "2014-11-20T03:24:10Z"

			Expect(docs.Diff(left, right)).To(BeFalse())
		})
	})
})
