package notifications_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientRegistrationParams", func() {
	Describe("NewClientRegistrationParams", func() {
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
			Expect(err).NotTo(HaveOccurred())

			parameters, err := notifications.NewClientRegistrationParams(bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			Expect(parameters.SourceName).To(Equal("Raptor Containment Unit"))
			Expect(len(parameters.Notifications)).To(Equal(2))
			Expect(parameters.Notifications).To(ContainElement(&notifications.NotificationStruct{
				ID:          "perimeter_breach",
				Description: "Perimeter Breach",
				Critical:    true,
			}))
			Expect(parameters.Notifications).To(ContainElement(&notifications.NotificationStruct{
				ID:          "feeding_time",
				Description: "Feeding Time",
				Critical:    false,
			}))
		})

		Context("error cases", func() {
			It("returns an error when the parameters are invalid JSON", func() {
				_, err := notifications.NewClientRegistrationParams(strings.NewReader("this is not valid JSON"))
				Expect(err).To(BeAssignableToTypeOf(webutil.ParseError{}))
			})

			It("returns an error when the request body is missing", func() {
				_, err := notifications.NewClientRegistrationParams(ErrorReader{})
				Expect(err).To(BeAssignableToTypeOf(webutil.ParseError{}))
			})

			Context("when the JSON contains unexpected properties", func() {
				It("returns an error for invalid top level keys ", func() {
					someJson := `{ "source_name" : "Raptor Containment Unit", "invalid_property" : 5 }`

					_, err := notifications.NewClientRegistrationParams(strings.NewReader(someJson))
					Expect(err).To(MatchError(webutil.SchemaError{Err: errors.New("\"invalid_property\" is not a valid property")}))
				})

				It("returns an error for invalid nested keys", func() {
					someJson := `{ "source_name" : "Raptor", "notifications": { "some_id": {"description" : "ok", "invalid_property" : 5 } } }`

					_, err := notifications.NewClientRegistrationParams(strings.NewReader(someJson))
					Expect(err).To(MatchError(webutil.SchemaError{Err: errors.New("\"invalid_property\" is not a valid property")}))
				})
			})

			Context("when the JSON contains null values", func() {
				It("returns an error if 'notifications' collection is null", func() {
					someJson := `{ "source_name" : "Something something raptor", "notifications": null }`

					_, err := notifications.NewClientRegistrationParams(strings.NewReader(someJson))
					Expect(err).To(MatchError(webutil.SchemaError{Err: errors.New("only include \"notifications\" key when adding a notification")}))
				})
				It("returns an error if an individual notification is null ", func() {
					someJson := `{ "source_name" : "Raptor", "notifications": { "some_id":  null } }`

					_, err := notifications.NewClientRegistrationParams(strings.NewReader(someJson))
					Expect(err).To(MatchError(webutil.SchemaError{Err: errors.New("notification must not be null")}))
				})
			})
		})
	})

	Describe("Validate", func() {
		It("validates when only source_name is present", func() {
			cr := notifications.ClientRegistrationParams{SourceName: "jurassic_park"}
			err := cr.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error when source_name is missing", func() {
			cr := notifications.ClientRegistrationParams{}
			err := cr.Validate()

			Expect(err).To(MatchError(webutil.ValidationError{Err: errors.New("\"source_name\" is a required field")}))
		})

		It("returns an error if notification is missing a required field", func() {
			cr := notifications.ClientRegistrationParams{
				SourceName: "jurassic_park",
				Notifications: map[string](*notifications.NotificationStruct){
					"perimeter_breach": {},
				},
			}

			err := cr.Validate()
			Expect(err).To(MatchError(webutil.ValidationError{
				Err: errors.New("notification \"perimeter_breach\" is missing required field \"ID\", notification \"perimeter_breach\" is missing required field \"Description\""),
			}))
		})
	})
})
