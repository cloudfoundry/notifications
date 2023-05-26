package webutil_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/v1/collections"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrorWriter", func() {
	var (
		writer   webutil.ErrorWriter
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		writer = webutil.NewErrorWriter()
		recorder = httptest.NewRecorder()
	})

	It("returns a 422 when a client tries to register a critical notification without critical_notifications.write scope", func() {
		writer.Write(recorder, webutil.UAAScopesError{Err: errors.New("UAA Scopes Error: Client does not have authority to register critical notifications.")})
		Expect(recorder.Code).To(Equal(422))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["UAA Scopes Error: Client does not have authority to register critical notifications."]
		}`))
	})

	It("returns a 502 when CloudController fails to respond", func() {
		writer.Write(recorder, services.CCDownError{Err: errors.New("Bad things happened!")})
		Expect(recorder.Code).To(Equal(http.StatusBadGateway))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["Bad things happened!"]
		}`))
	})

	It("returns a 500 when there is a template create error", func() {
		writer.Write(recorder, webutil.TemplateCreateError{})
		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["Failed to create Template in the database"]
		}`))
	})

	It("returns a 500 when there is a template update error", func() {
		writer.Write(recorder, models.TemplateUpdateError{Err: errors.New("Failed to update Template in the database")})
		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["Failed to update Template in the database"]
		}`))
	})

	It("returns a 404 when the space cannot be found", func() {
		writer.Write(recorder, services.CCNotFoundError{Err: errors.New("Space could not be found")})
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["Space could not be found"]
		}`))
	})

	It("returns a 404 when the space cannot be found", func() {
		writer.Write(recorder, cf.NotFoundError{Message: "Space could not be found"})
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["CloudController Failure: Space could not be found"]
		}`))
	})

	It("returns a 400 when the request cannot be parsed due to syntatically invalid JSON", func() {
		writer.Write(recorder, webutil.ParseError{})
		Expect(recorder.Code).To(Equal(400))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["Request body could not be parsed"]
		}`))
	})

	It("returns a 422 when the requests are not valid due to semantically invalid JSON", func() {
		writer.Write(recorder, webutil.ValidationError{Err: errors.New("invalid json")})
		Expect(recorder.Code).To(Equal(422))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["invalid json"]
		}`))
	})

	It("returns a 422 when trying to send a critical notification without correct scope", func() {
		writer.Write(recorder, webutil.NewCriticalNotificationError("raptors"))
		Expect(recorder.Code).To(Equal(422))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["Insufficient privileges to send notification raptors"]
		}`))
	})

	It("returns a 409 when there is a duplicate record", func() {
		writer.Write(recorder, models.DuplicateError{Err: errors.New("duplicate record")})
		Expect(recorder.Code).To(Equal(409))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["duplicate record"]
		}`))
	})

	It("returns a 404 when a record cannot be found", func() {
		writer.Write(recorder, models.NotFoundError{Err: errors.New("not found")})
		Expect(recorder.Code).To(Equal(404))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["not found"]
		}`))
	})

	It("returns a 406 when a record cannot be found", func() {
		writer.Write(recorder, services.DefaultScopeError{})
		Expect(recorder.Code).To(Equal(406))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["You cannot send a notification to a default scope"]
		}`))
	})

	It("returns a 422 when a template cannot be assigned", func() {
		writer.Write(recorder, collections.TemplateAssignmentError{Err: errors.New("The template could not be assigned")})
		Expect(recorder.Code).To(Equal(422))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["The template could not be assigned"]
		}`))
	})

	It("returns a 422 when a user token was expected but is not present", func() {
		writer.Write(recorder, webutil.MissingUserTokenError{Err: errors.New("Missing user_id from token claims.")})
		Expect(recorder.Code).To(Equal(422))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["Missing user_id from token claims."]
		}`))
	})

	It("returns a 500 for unknown errors", func() {
		writer.Write(recorder, errors.New("unknown error"))
		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).To(MatchJSON(`{
			"errors": ["unknown error"]
		}`))
	})
})
