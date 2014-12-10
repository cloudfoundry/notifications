package valiant_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/valiant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	Context("when the data is one level deep", func() {
		It("succeeds when the json is valid", func() {
			data := `{"name":"Boshy", "email": true}`

			type Person struct {
				Name  string `json:"name"    validate-required:"true"`
				Email bool   `json:"email"   validate-required:"false"`
			}

			var someone Person

			err := json.Unmarshal([]byte(data), &someone)
			if err != nil {
				panic(err)
			}

			err = valiant.ValidateJSON(someone, []byte(data))
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error when the json is missing a required field", func() {
			data := `{"email": true}`

			type Person struct {
				Name  string `json:"name"    validate-required:"true"`
				Email bool   `json:"email"   validate-required:"false"`
			}

			var someone Person
			err := json.Unmarshal([]byte(data), &someone)
			if err != nil {
				panic(err)
			}

			err = valiant.ValidateJSON(someone, []byte(data))
			Expect(err).To(HaveOccurred())

		})
	})

	Context("when the data is nested n levels deep", func() {
		It("succeeds when the json is valid, even if validate-required tag is not present", func() {
			data := `{
				"name":"Boshy",
				"contact_info": {
				    "address": {
					    "street":"123 Sesame St",
					    "city":"Santa Monica",
					    "state":"CA"
				    },
				"phone": "310-310-3100",
				"email": true
			}}`

			type Address struct {
				Street string `json:"street" validate-required:"false"`
				City   string `json:"city" validate-required:"true"`
				State  string `json:"state"`
			}

			type ContactInfo struct {
				Address Address `json:"address" validate-required:"false"`
				Phone   string  `json:"phone" validate-required:"true"`
				Email   bool    `json:"email"   validate-required:"false"`
			}

			type Person struct {
				Name        string      `json:"name"    validate-required:"true"`
				ContactInfo ContactInfo `json:"contact_info" validate-required:"true"`
			}

			var someone Person

			err := json.Unmarshal([]byte(data), &someone)
			if err != nil {
				panic(err)
			}

			err = valiant.ValidateJSON(someone, []byte(data))
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error when the json is missing a required field", func() {
			data := `{
				"name":"Boshy",
				"contact_info": {
				    "address": {
					    "street":"123 Sesame St"
				    },
				"phone": "310-310-3100",
				"email": true
			}}`

			type Address struct {
				Street string `json:"street" validate-required:"false"`
				City   string `json:"city" validate-required:"true"`
				State  string `json:"state" validate-required:"true"`
			}

			type ContactInfo struct {
				Address Address `json:"address" validate-required:"false"`
				Phone   string  `json:"phone" validate-required:"true"`
				Email   bool    `json:"email"   validate-required:"false"`
			}

			type Person struct {
				Name        string      `json:"name"    validate-required:"true"`
				ContactInfo ContactInfo `json:"contact_info" validate-required:"true"`
			}

			var someone Person

			err := json.Unmarshal([]byte(data), &someone)
			if err != nil {
				panic(err)
			}

			err = valiant.ValidateJSON(someone, []byte(data))
			Expect(err).To(HaveOccurred())

		})
	})
})
