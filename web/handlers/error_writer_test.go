package handlers_test

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("ErrorWriter", func() {
    var writer handlers.ErrorWriter
    var recorder *httptest.ResponseRecorder

    BeforeEach(func() {
        writer = handlers.NewErrorWriter()
        recorder = httptest.NewRecorder()
    })

    It("returns a 502 when CloudController fails to respond", func() {
        writer.Write(recorder, postal.CCDownError("BOOM!"))

        Expect(recorder.Code).To(Equal(http.StatusBadGateway))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("Cloud Controller is unavailable"))
    })

    It("returns a 502 when UAA fails to respond", func() {
        writer.Write(recorder, postal.UAADownError("BOOM!"))

        Expect(recorder.Code).To(Equal(http.StatusBadGateway))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("UAA is unavailable"))
    })

    It("returns a 502 when UAA fails for unknown reasons", func() {
        writer.Write(recorder, postal.UAAGenericError("UAA Unknown Error: BOOM!"))

        Expect(recorder.Code).To(Equal(http.StatusBadGateway))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("UAA Unknown Error: BOOM!"))
    })

    It("returns a 500 when the is a template loading error", func() {
        writer.Write(recorder, postal.TemplateLoadError("BOOM!"))

        Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("An email template could not be loaded"))
    })

    It("returns a 404 when the space cannot be found", func() {
        writer.Write(recorder, postal.CCNotFoundError("Organization could not be found"))

        Expect(recorder.Code).To(Equal(http.StatusNotFound))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("CloudController Error: Organization could not be found"))
    })

    It("returns a 422 when the params cannot be parsed", func() {
        writer.Write(recorder, handlers.ParamsParseError{})

        Expect(recorder.Code).To(Equal(422))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("Request body could not be parsed"))
    })

    It("returns a 422 when the params are not valid", func() {
        writer.Write(recorder, handlers.ParamsValidationError([]string{"something", "another"}))

        Expect(recorder.Code).To(Equal(422))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("something"))
        Expect(body["errors"]).To(ContainElement("another"))
    })

    It("returns a 409 when there is a duplicate record", func() {
        writer.Write(recorder, models.ErrDuplicateRecord{})

        Expect(recorder.Code).To(Equal(409))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("Duplicate Record"))
    })

    It("returns a 404 when a record cannot be found", func() {
        writer.Write(recorder, models.ErrRecordNotFound{})

        Expect(recorder.Code).To(Equal(404))

        body := make(map[string]interface{})
        err := json.Unmarshal(recorder.Body.Bytes(), &body)
        if err != nil {
            panic(err)
        }

        Expect(body["errors"]).To(ContainElement("Record Not Found"))
    })

    It("panics for unknown errors", func() {
        Expect(func() {
            writer.Write(recorder, errors.New("BOOM!"))
        }).To(Panic())
    })
})
