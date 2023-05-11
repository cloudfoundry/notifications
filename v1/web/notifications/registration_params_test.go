package notifications_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ErrorReader struct{}

func (reader ErrorReader) Read(b []byte) (int, error) {
	return 0, errors.New("BOOM!")
}

func (reader ErrorReader) Close() error {
	return nil
}

var _ = Describe("RegistrationParams", func() {
	Describe("NewRegistrationParams", func() {
		It("constructs parameters from a reader", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_description": "Raptor Containment Unit",
				"kinds": []map[string]interface{}{
					{
						"id":          "perimeter_breach",
						"description": "Perimeter Breach",
						"critical":    true,
					},
					{
						"id":          "feeding_time",
						"description": "Feeding Time",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			parameters, err := notifications.NewRegistrationParams(ioutil.NopCloser(bytes.NewBuffer(body)))
			Expect(err).NotTo(HaveOccurred())
			Expect(parameters.SourceDescription).To(Equal("Raptor Containment Unit"))
			Expect(len(parameters.Kinds)).To(Equal(2))
			Expect(parameters.Kinds).To(ContainElement(models.Kind{
				ID:          "perimeter_breach",
				Description: "Perimeter Breach",
				Critical:    true,
			}))
			Expect(parameters.Kinds).To(ContainElement(models.Kind{
				ID:          "feeding_time",
				Description: "Feeding Time",
				Critical:    false,
			}))
			Expect(parameters.IncludesKinds).To(BeTrue())
		})

		It("sets the IncludesKinds flag to false when the kinds are missing", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_description": "Raptor Containment Unit",
			})
			Expect(err).NotTo(HaveOccurred())

			parameters, err := notifications.NewRegistrationParams(ioutil.NopCloser(bytes.NewBuffer(body)))
			Expect(err).NotTo(HaveOccurred())

			Expect(parameters.IncludesKinds).To(BeFalse())
		})

		It("returns an error when the parameters are invalid JSON", func() {
			_, err := notifications.NewRegistrationParams(ioutil.NopCloser(strings.NewReader("this is not valid JSON")))
			Expect(err).To(BeAssignableToTypeOf(webutil.ParseError{}))
		})

		It("returns an error when the request body is missing", func() {
			_, err := notifications.NewRegistrationParams(ErrorReader{})
			Expect(err).To(BeAssignableToTypeOf(webutil.ParseError{}))
		})
	})

	Describe("Validate", func() {
		It("validates a valid request body", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_description": "Raptor Containment Unit",
				"kinds": []map[string]interface{}{
					{
						"id":          "perimeter_breach-88._",
						"description": "Perimeter Breach",
						"critical":    true,
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			parameters, err := notifications.NewRegistrationParams(ioutil.NopCloser(bytes.NewBuffer(body)))
			Expect(err).NotTo(HaveOccurred())

			err = parameters.Validate()
			Expect(err).To(BeNil())
		})

		It("validates the presence of source_description, kind.id and kind.description", func() {
			body, err := json.Marshal(map[string]interface{}{
				"kinds": []models.Kind{
					{Critical: false},
					{Critical: false},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			parameters, err := notifications.NewRegistrationParams(ioutil.NopCloser(bytes.NewBuffer(body)))
			Expect(err).NotTo(HaveOccurred())

			err = parameters.Validate()
			Expect(err).To(MatchError(webutil.ValidationError{Err: errors.New("\"source_description\" is a required field, \"kind.id\" is a required field, \"kind.description\" is a required field")}))
		})

		It("validates the format of kind.ID's", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_description": "the source description",
				"kinds": []models.Kind{
					{
						ID:          "not-Valid@",
						Description: "kind description",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			parameters, err := notifications.NewRegistrationParams(ioutil.NopCloser(bytes.NewBuffer(body)))
			Expect(err).NotTo(HaveOccurred())

			err = parameters.Validate()
			Expect(err).To(MatchError(webutil.ValidationError{Err: errors.New("\"kind.id\" is improperly formatted")}))
		})

	})
})
