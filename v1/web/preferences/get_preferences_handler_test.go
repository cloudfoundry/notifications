package preferences_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetPreferencesHandler", func() {
	var (
		handler           preferences.GetPreferencesHandler
		writer            *httptest.ResponseRecorder
		request           *http.Request
		preferencesFinder *mocks.PreferencesFinder
		errorWriter       *mocks.ErrorWriter
		builder           services.PreferencesBuilder
		context           stack.Context
		database          *mocks.Database

		TRUE  = true
		FALSE = false
	)

	BeforeEach(func() {
		errorWriter = mocks.NewErrorWriter()

		writer = httptest.NewRecorder()
		body, err := json.Marshal(map[string]string{
			"I think this request is empty": "maybe",
		})
		if err != nil {
			panic(err)
		}

		tokenHeader := map[string]interface{}{
			"alg": "RS256",
		}
		tokenClaims := map[string]interface{}{
			"user_id": "correct-user",
			"exp":     int64(3404281214),
			"scope":   []string{"notification_preferences.read"},
		}

		request, err = http.NewRequest("GET", "/user_preferences", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		token, err := jwt.Parse(helpers.BuildToken(tokenHeader, tokenClaims), func(token *jwt.Token) (interface{}, error) {
			return helpers.UAAPublicKeyRSA, nil
		})

		database = mocks.NewDatabase()

		context = stack.NewContext()
		context.Set("token", token)
		context.Set("database", database)

		builder = services.NewPreferencesBuilder()
		builder.Add(models.Preference{
			ClientID: "raptorClient",
			KindID:   "hungry-kind",
			Email:    false,
		})
		builder.Add(models.Preference{
			ClientID: "starWarsClient",
			KindID:   "vader-kind",
			Email:    true,
		})
		builder.GlobalUnsubscribe = true

		preferencesFinder = mocks.NewPreferencesFinder()
		preferencesFinder.FindCall.Returns.PreferencesBuilder = builder

		handler = preferences.NewGetPreferencesHandler(preferencesFinder, errorWriter)
	})

	It("passes the proper user guid into execute", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(preferencesFinder.FindCall.Receives.Database).To(Equal(database))
		Expect(preferencesFinder.FindCall.Receives.UserGUID).To(Equal("correct-user"))
	})

	It("returns a proper JSON response when the Preference object does not error", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))

		parsed := services.PreferencesBuilder{}
		err := json.Unmarshal(writer.Body.Bytes(), &parsed)
		if err != nil {
			panic(err)
		}

		Expect(parsed.GlobalUnsubscribe).To(BeTrue())
		Expect(parsed.Clients["raptorClient"]["hungry-kind"].Email).To(Equal(&FALSE))
		Expect(parsed.Clients["starWarsClient"]["vader-kind"].Email).To(Equal(&TRUE))
	})

	Context("when there is an error returned from the finder", func() {
		It("writes the error to the error writer", func() {
			preferencesFinder.FindCall.Returns.Error = errors.New("boom!")
			handler.ServeHTTP(writer, request, context)
			Expect(errorWriter.WriteCall.Receives.Error).To(Equal(preferencesFinder.FindCall.Returns.Error))
		})
	})

	Context("when the request does not container a user_id claim", func() {
		It("writes a MissingUserTokenError to the error writer", func() {
			tokenHeader := map[string]interface{}{
				"alg": "RS256",
			}

			tokenClaims := map[string]interface{}{}

			request, err := http.NewRequest("GET", "/user_preferences", nil)
			if err != nil {
				panic(err)
			}

			token, err := jwt.Parse(helpers.BuildToken(tokenHeader, tokenClaims), func(token *jwt.Token) (interface{}, error) {
				return helpers.UAAPublicKeyRSA, nil
			})

			context = stack.NewContext()
			context.Set("token", token)

			handler.ServeHTTP(writer, request, context)
			Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.MissingUserTokenError{Err: errors.New("Missing user_id from token claims.")}))
		})
	})
})
