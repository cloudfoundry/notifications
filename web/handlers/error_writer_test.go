package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"

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

	It("returns a 422 when a client tries to register a critical notification without critical_notifications.write scope", func() {
		writer.Write(recorder, postal.UAAScopesError("UAA Scopes Error: Client does not have authority to register critical notifications."))

		unprocessableEntity := 422
		Expect(recorder.Code).To(Equal(unprocessableEntity))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("UAA Scopes Error: Client does not have authority to register critical notifications."))
	})

	It("returns a 502 when CloudController fails to respond", func() {
		writer.Write(recorder, utilities.CCDownError("Bad things happened!"))

		Expect(recorder.Code).To(Equal(http.StatusBadGateway))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Bad things happened!"))
	})

	It("returns a 502 when UAA fails to respond", func() {
		writer.Write(recorder, utilities.UAADownError("Whoops!"))

		Expect(recorder.Code).To(Equal(http.StatusBadGateway))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Whoops!"))
	})

	It("returns a 502 when UAA fails for unknown reasons", func() {
		writer.Write(recorder, utilities.UAAGenericError("UAA Unknown Error: BOOM!"))

		Expect(recorder.Code).To(Equal(http.StatusBadGateway))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("UAA Unknown Error: BOOM!"))
	})

	It("returns a 500 when there is a template loading error", func() {
		writer.Write(recorder, postal.TemplateLoadError("BOOM!"))

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("An email template could not be loaded"))
	})

	It("returns a 500 when there is a template create error", func() {
		writer.Write(recorder, params.TemplateCreateError{})

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Failed to create Template in the database"))
	})

	It("returns a 404 when there is a template find error", func() {
		writer.Write(recorder, models.TemplateFindError{Message: "Template my-id could not be found"})

		Expect(recorder.Code).To(Equal(http.StatusNotFound))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Template my-id could not be found"))
	})

	It("returns a 500 when there is a template update error", func() {
		writer.Write(recorder, models.TemplateUpdateError{Message: "Failed to update Template in the database"})

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Failed to update Template in the database"))
	})

	It("returns a 404 when the space cannot be found", func() {
		writer.Write(recorder, utilities.CCNotFoundError("Organization could not be found"))

		Expect(recorder.Code).To(Equal(http.StatusNotFound))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("CloudController Error: Organization could not be found"))
	})

	It("returns a 400 when the params cannot be parsed due to syntatically invalid JSON", func() {
		writer.Write(recorder, params.ParseError{})

		Expect(recorder.Code).To(Equal(400))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Request body could not be parsed"))
	})

	It("returns a 422 when the params are not valid due to semantically invalid JSON", func() {
		writer.Write(recorder, params.ValidationError([]string{"something", "another"}))

		Expect(recorder.Code).To(Equal(422))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("something"))
		Expect(body["errors"]).To(ContainElement("another"))
	})

	It("returns a 422 when trying to send a critical notification without correct scope", func() {
		writer.Write(recorder, postal.NewCriticalNotificationError("raptors"))

		Expect(recorder.Code).To(Equal(422))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Insufficient privileges to send notification raptors"))
	})

	It("returns a 409 when there is a duplicate record", func() {
		writer.Write(recorder, models.DuplicateRecordError{})

		Expect(recorder.Code).To(Equal(409))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Duplicate Record"))
	})

	It("returns a 404 when a record cannot be found", func() {
		writer.Write(recorder, models.NewRecordNotFoundError("hello"))

		Expect(recorder.Code).To(Equal(404))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Record Not Found: hello"))
	})

	It("returns a 406 when a record cannot be found", func() {
		writer.Write(recorder, strategies.DefaultScopeError{})
		Expect(recorder.Code).To(Equal(406))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("You cannot send a notification to a default scope"))
	})

	It("returns a 422 when a template cannot be assigned", func() {
		writer.Write(recorder, services.TemplateAssignmentError("The template could not be assigned"))
		Expect(recorder.Code).To(Equal(422))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("The template could not be assigned"))
	})

	It("returns a 422 when a user token was expected but is not present", func() {
		writer.Write(recorder, handlers.MissingUserTokenError("Missing user_id from token claims."))
		Expect(recorder.Code).To(Equal(422))

		body := make(map[string]interface{})
		err := json.Unmarshal(recorder.Body.Bytes(), &body)
		if err != nil {
			panic(err)
		}

		Expect(body["errors"]).To(ContainElement("Missing user_id from token claims."))
	})

	It("panics for unknown errors", func() {
		Expect(func() {
			writer.Write(recorder, errors.New("BOOM!"))
		}).To(Panic())
	})
})
