package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

type ErrorReader struct{}

func (reader ErrorReader) Read(b []byte) (int, error) {
    return 0, errors.New("BOOM!")
}

var _ = Describe("RegistrationParams", func() {
    Describe("NewRegistrationParams", func() {
        It("constructs params from a reader", func() {
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
            if err != nil {
                panic(err)
            }

            params, err := handlers.NewRegistrationParams(bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            Expect(params.SourceDescription).To(Equal("Raptor Containment Unit"))
            Expect(len(params.Kinds)).To(Equal(2))
            Expect(params.Kinds).To(ContainElement(models.Kind{
                ID:          "perimeter_breach",
                Description: "Perimeter Breach",
                Critical:    true,
            }))
            Expect(params.Kinds).To(ContainElement(models.Kind{
                ID:          "feeding_time",
                Description: "Feeding Time",
                Critical:    false,
            }))
        })

        It("returns an error when the params are invalid JSON", func() {
            _, err := handlers.NewRegistrationParams(strings.NewReader("this is not valid JSON"))
            Expect(err).To(BeAssignableToTypeOf(handlers.ParamsParseError{}))
        })

        It("returns an error when the request body is missing", func() {
            _, err := handlers.NewRegistrationParams(ErrorReader{})
            Expect(err).To(BeAssignableToTypeOf(handlers.ParamsParseError{}))
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
            if err != nil {
                panic(err)
            }

            params, err := handlers.NewRegistrationParams(bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            err = params.Validate()
            Expect(err).To(BeNil())
        })

        It("validates the presence of source_description, kind.id and kind.description", func() {
            body, err := json.Marshal(map[string]interface{}{
                "kinds": []models.Kind{
                    {Critical: false},
                    {Critical: false},
                },
            })
            if err != nil {
                panic(err)
            }

            params, err := handlers.NewRegistrationParams(bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            err = params.Validate()
            Expect(err).To(BeAssignableToTypeOf(handlers.ParamsValidationError{}))
            errs := err.(handlers.ParamsValidationError).Errors()
            Expect(len(errs)).To(Equal(3))
            Expect(err).To(ContainElement(`"source_description" is a required field`))
            Expect(err).To(ContainElement(`"kind.id" is a required field`))
            Expect(err).To(ContainElement(`"kind.description" is a required field`))
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
            if err != nil {
                panic(err)
            }

            params, err := handlers.NewRegistrationParams(bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }
            err = params.Validate()
            Expect(err).To(BeAssignableToTypeOf(handlers.ParamsValidationError{}))
            errs := err.(handlers.ParamsValidationError).Errors()
            Expect(len(errs)).To(Equal(1))
            Expect(err).To(ContainElement(`"kind.id" is improperly formatted`))
        })

    })
})
