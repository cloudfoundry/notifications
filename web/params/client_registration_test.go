package params_test

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientRegistration", func() {
	Describe("NewClientRegistration", func() {
		It("constructs parameters from a reader", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_name": "Raptor Containment Unit",
				"notifications": map[string]interface{}{
					"perimeter_breach": map[string]interface{}{
						"description": "Perimeter Breach",
						"critical":    true,
					},
					"feeding_time": map[string]interface{}{
						"description": "Feeding Time",
					},
				},
			})
			if err != nil {
				panic(err)
			}

			parameters, err := params.NewClientRegistration(bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}

			Expect(parameters.SourceName).To(Equal("Raptor Containment Unit"))
			Expect(len(parameters.Notifications)).To(Equal(2))
			Expect(parameters.Notifications).To(ContainElement(&params.NotificationStruct{
				ID:          "perimeter_breach",
				Description: "Perimeter Breach",
				Critical:    true,
			}))
			Expect(parameters.Notifications).To(ContainElement(&params.NotificationStruct{
				ID:          "feeding_time",
				Description: "Feeding Time",
				Critical:    false,
			}))
		})

		Context("error cases", func() {
			It("returns an error when the parameters are invalid JSON", func() {
				_, err := params.NewClientRegistration(strings.NewReader("this is not valid JSON"))
				Expect(err).To(BeAssignableToTypeOf(params.ParseError{}))
			})

			It("returns an error when the request body is missing", func() {
				_, err := params.NewClientRegistration(ErrorReader{})
				Expect(err).To(BeAssignableToTypeOf(params.ParseError{}))
			})

			Context("when the JSON contains unexpected properties", func() {
				It("returns an error for invalid top level keys ", func() {
					someJson := `{ "source_name" : "Raptor Containment Unit", "invalid_property" : 5 }`
					_, err := params.NewClientRegistration(strings.NewReader(someJson))
					Expect(err).To(BeAssignableToTypeOf(params.NewSchemaError("")))
				})

				It("returns an error for invalid nested keys", func() {
					someJson := `{ "source_name" : "Raptor", "notifications": { "some_id": {"description" : "ok", "invalid_property" : 5 } } }`
					_, err := params.NewClientRegistration(strings.NewReader(someJson))
					Expect(err).To(BeAssignableToTypeOf(params.NewSchemaError("")))
				})
			})

			Context("when the JSON contains null values", func() {
				It("returns an error if 'notifications' collection is null", func() {
					someJson := `{ "source_name" : "Something something raptor", "notifications": null }`
					_, err := params.NewClientRegistration(strings.NewReader(someJson))
					Expect(err).To(BeAssignableToTypeOf(params.NewSchemaError("")))
				})
				It("returns an error if an individual notification is null ", func() {
					someJson := `{ "source_name" : "Raptor", "notifications": { "some_id":  null } }`
					_, err := params.NewClientRegistration(strings.NewReader(someJson))
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(params.NewSchemaError("")))
				})
			})
		})
	})

	Describe("Validate", func() {
		It("validates when only source_name is present", func() {
			cr := params.ClientRegistration{SourceName: "jurassic_park"}
			err := cr.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error when source_name is missing", func() {
			cr := params.ClientRegistration{}
			err := cr.Validate()

			Expect(err).To(BeAssignableToTypeOf(params.ValidationError{}))
			errs := err.(params.ValidationError).Errors()
			Expect(len(errs)).To(Equal(1))
			Expect(err).To(ContainElement(ContainSubstring("source_name")))
		})

		It("returns an error if notification is missing a required field", func() {
			cr := params.ClientRegistration{
				SourceName: "jurassic_park",
				Notifications: map[string](*params.NotificationStruct){
					"perimeter_breach": {},
				},
			}
			err := cr.Validate()

			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(params.ValidationError{}))
			Expect(err).To(ContainElement(`notification "perimeter_breach" is missing required field "ID"`))
			Expect(err).To(ContainElement(`notification "perimeter_breach" is missing required field "Description"`))
		})

	})
})
