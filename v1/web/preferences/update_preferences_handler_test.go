package preferences_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

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

var _ = Describe("UpdatePreferencesHandler", func() {
	Describe("Execute", func() {
		var (
			handler     preferences.UpdatePreferencesHandler
			writer      *httptest.ResponseRecorder
			request     *http.Request
			updater     *mocks.PreferenceUpdater
			errorWriter *mocks.ErrorWriter
			conn        *mocks.Connection
			transaction *mocks.Transaction
			context     stack.Context
		)

		BeforeEach(func() {
			transaction = mocks.NewTransaction()

			conn = mocks.NewConnection()
			conn.TransactionCall.Returns.Transaction = transaction

			database := mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = conn

			builder := services.NewPreferencesBuilder()

			builder.Add(models.Preference{
				ClientID: "raptors",
				KindID:   "door-opening",
				Email:    false,
			})
			builder.Add(models.Preference{
				ClientID: "raptors",
				KindID:   "feeding-time",
				Email:    true,
			})
			builder.Add(models.Preference{
				ClientID: "dogs",
				KindID:   "barking",
				Email:    false,
			})
			builder.GlobalUnsubscribe = true

			body, err := json.Marshal(builder)
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			tokenHeader := map[string]interface{}{
				"alg": "RS256",
			}
			tokenClaims := map[string]interface{}{
				"user_id": "correct-user",
				"exp":     int64(3404281214),
			}

			rawToken := helpers.BuildToken(tokenHeader, tokenClaims)
			request.Header.Set("Authorization", "Bearer "+rawToken)
			token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
				return helpers.UAAPublicKeyRSA, nil
			})
			Expect(err).NotTo(HaveOccurred())

			context = stack.NewContext()
			context.Set("token", token)
			context.Set("database", database)

			errorWriter = mocks.NewErrorWriter()
			updater = mocks.NewPreferenceUpdater()
			writer = httptest.NewRecorder()

			handler = preferences.NewUpdatePreferencesHandler(updater, errorWriter)
		})

		It("Passes The Correct Arguments to PreferenceUpdater Execute", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(reflect.ValueOf(updater.UpdateCall.Receives.Connection).Pointer()).To(Equal(reflect.ValueOf(transaction).Pointer()))

			preferencesArguments := updater.UpdateCall.Receives.Preferences

			Expect(preferencesArguments).To(ContainElement(models.Preference{
				ClientID: "raptors",
				KindID:   "door-opening",
				Email:    false,
			}))
			Expect(preferencesArguments).To(ContainElement(models.Preference{
				ClientID: "raptors",
				KindID:   "feeding-time",
				Email:    true,
			}))
			Expect(preferencesArguments).To(ContainElement(models.Preference{
				ClientID: "dogs",
				KindID:   "barking",
				Email:    false,
			}))

			Expect(updater.UpdateCall.Receives.GlobalUnsubscribe).To(BeTrue())
			Expect(updater.UpdateCall.Receives.UserID).To(Equal("correct-user"))
		})

		It("Returns a 204 status code when the Preference object does not error", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("Failure cases", func() {
			It("returns an error when the clients key is missing", func() {
				jsonBody := `{"raptor-client": {"containment-unit-breach": {"email": false}}}`

				request, err := http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer([]byte(jsonBody)))
				Expect(err).NotTo(HaveOccurred())

				tokenHeader := map[string]interface{}{
					"alg": "RS256",
				}
				tokenClaims := map[string]interface{}{
					"user_id": "correct-user",
					"exp":     int64(3404281214),
				}

				rawToken := helpers.BuildToken(tokenHeader, tokenClaims)
				request.Header.Set("Authorization", "Bearer "+rawToken)

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).ToNot(BeNil())
				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))

				Expect(transaction.BeginCall.WasCalled).To(BeFalse())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
			})

			Context("preferenceUpdater.Execute errors", func() {
				Context("when the user_id claim is not present in the token", func() {
					It("Writes a MissingUserTokenError to the error writer", func() {
						tokenHeader := map[string]interface{}{
							"alg": "RS256",
						}

						tokenClaims := map[string]interface{}{}

						request, err := http.NewRequest("PATCH", "/user_preferences", nil)
						Expect(err).NotTo(HaveOccurred())

						token, err := jwt.Parse(helpers.BuildToken(tokenHeader, tokenClaims), func(token *jwt.Token) (interface{}, error) {
							return helpers.UAAPublicKeyRSA, nil
						})
						Expect(err).NotTo(HaveOccurred())

						context.Set("token", token)

						handler.ServeHTTP(writer, request, context)
						Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.MissingUserTokenError{Err: errors.New("Missing user_id from token claims.")}))
						Expect(transaction.BeginCall.WasCalled).To(BeFalse())
						Expect(transaction.CommitCall.WasCalled).To(BeFalse())
						Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
					})
				})

				It("delegates MissingKindOrClientErrors as webutil.ValidationError to the ErrorWriter", func() {
					updateError := services.MissingKindOrClientError{Err: errors.New("BOOM!")}
					updater.UpdateCall.Returns.Error = updateError

					handler.ServeHTTP(writer, request, context)

					Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.ValidationError{Err: updateError}))

					Expect(transaction.BeginCall.WasCalled).To(BeTrue())
					Expect(transaction.CommitCall.WasCalled).To(BeFalse())
					Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
				})

				It("delegates CriticalKindErrors as webutil.ValidationError to the ErrorWriter", func() {
					updateError := services.CriticalKindError{Err: errors.New("BOOM!")}
					updater.UpdateCall.Returns.Error = updateError

					handler.ServeHTTP(writer, request, context)

					Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.ValidationError{Err: updateError}))

					Expect(transaction.BeginCall.WasCalled).To(BeTrue())
					Expect(transaction.CommitCall.WasCalled).To(BeFalse())
					Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
				})

				It("delegates other errors to the ErrorWriter", func() {
					updater.UpdateCall.Returns.Error = errors.New("BOOM!")

					handler.ServeHTTP(writer, request, context)

					Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BOOM!")))

					Expect(transaction.BeginCall.WasCalled).To(BeTrue())
					Expect(transaction.CommitCall.WasCalled).To(BeFalse())
					Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
				})
			})

			It("delegates json validation errors to the ErrorWriter", func() {
				requestBody, err := json.Marshal(map[string]interface{}{
					"something": true,
				})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))

				Expect(transaction.BeginCall.WasCalled).To(BeFalse())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
			})

			It("delegates validation errors to the error writer", func() {
				requestBody, err := json.Marshal(map[string]map[string]map[string]map[string]interface{}{
					"clients": {
						"client-id": {
							"kind-id": {},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))

				Expect(transaction.BeginCall.WasCalled).To(BeFalse())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
			})

			It("delegates transaction errors to the error writer", func() {
				transaction.CommitCall.Returns.Error = errors.New("transaction error, oh no")
				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(models.TransactionCommitError{Err: errors.New("transaction error, oh no")}))

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeTrue())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
			})
		})
	})
})
