package preferences_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdatePreferencesHandler", func() {
	Describe("Execute", func() {
		var (
			handler     preferences.UpdatePreferencesHandler
			writer      *httptest.ResponseRecorder
			request     *http.Request
			updater     *fakes.PreferenceUpdater
			errorWriter *fakes.ErrorWriter
			conn        *fakes.Connection
			context     stack.Context
		)

		BeforeEach(func() {
			conn = fakes.NewConnection()
			database := fakes.NewDatabase()
			database.Conn = conn
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
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"user_id": "correct-user",
				"exp":     int64(3404281214),
			}

			rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
			request.Header.Set("Authorization", "Bearer "+rawToken)

			token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
				return []byte(application.UAAPublicKey), nil
			})
			Expect(err).NotTo(HaveOccurred())

			context = stack.NewContext()
			context.Set("token", token)
			context.Set("database", database)

			errorWriter = fakes.NewErrorWriter()
			updater = fakes.NewPreferenceUpdater()
			writer = httptest.NewRecorder()

			handler = preferences.NewUpdatePreferencesHandler(updater, errorWriter)
		})

		It("Passes The Correct Arguments to PreferenceUpdater Execute", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(len(updater.ExecuteArguments)).To(Equal(4))

			Expect(reflect.ValueOf(updater.ExecuteArguments[0]).Pointer()).To(Equal(reflect.ValueOf(conn).Pointer()))

			preferencesArguments := updater.ExecuteArguments[1]

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

			Expect(updater.ExecuteArguments[2]).To(BeTrue())
			Expect(updater.ExecuteArguments[3]).To(Equal("correct-user"))
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
					"alg": "FAST",
				}
				tokenClaims := map[string]interface{}{
					"user_id": "correct-user",
					"exp":     int64(3404281214),
				}

				rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
				request.Header.Set("Authorization", "Bearer "+rawToken)

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).ToNot(BeNil())
				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))
				Expect(conn.BeginWasCalled).To(BeFalse())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeFalse())
			})

			Context("preferenceUpdater.Execute errors", func() {
				Context("when the user_id claim is not present in the token", func() {
					It("Writes a MissingUserTokenError to the error writer", func() {
						tokenHeader := map[string]interface{}{
							"alg": "FAST",
						}

						tokenClaims := map[string]interface{}{}

						request, err := http.NewRequest("PATCH", "/user_preferences", nil)
						Expect(err).NotTo(HaveOccurred())

						token, err := jwt.Parse(fakes.BuildToken(tokenHeader, tokenClaims), func(token *jwt.Token) (interface{}, error) {
							return []byte(application.UAAPublicKey), nil
						})
						Expect(err).NotTo(HaveOccurred())

						context.Set("token", token)

						handler.ServeHTTP(writer, request, context)
						Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.MissingUserTokenError("")))
						Expect(conn.BeginWasCalled).To(BeFalse())
						Expect(conn.CommitWasCalled).To(BeFalse())
						Expect(conn.RollbackWasCalled).To(BeFalse())
					})
				})

				It("delegates MissingKindOrClientErrors as webutil.ValidationError to the ErrorWriter", func() {
					updater.ExecuteError = services.MissingKindOrClientError("BOOM!")

					handler.ServeHTTP(writer, request, context)

					Expect(errorWriter.WriteCall.Receives.Error).To(Equal(webutil.ValidationError([]string{"BOOM!"})))

					Expect(conn.BeginWasCalled).To(BeTrue())
					Expect(conn.CommitWasCalled).To(BeFalse())
					Expect(conn.RollbackWasCalled).To(BeTrue())
				})

				It("delegates CriticalKindErrors as webutil.ValidationError to the ErrorWriter", func() {
					updater.ExecuteError = services.CriticalKindError("BOOM!")

					handler.ServeHTTP(writer, request, context)

					Expect(errorWriter.WriteCall.Receives.Error).To(Equal(webutil.ValidationError([]string{"BOOM!"})))

					Expect(conn.BeginWasCalled).To(BeTrue())
					Expect(conn.CommitWasCalled).To(BeFalse())
					Expect(conn.RollbackWasCalled).To(BeTrue())
				})

				It("delegates other errors to the ErrorWriter", func() {
					updater.ExecuteError = errors.New("BOOM!")

					handler.ServeHTTP(writer, request, context)

					Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BOOM!")))

					Expect(conn.BeginWasCalled).To(BeTrue())
					Expect(conn.CommitWasCalled).To(BeFalse())
					Expect(conn.RollbackWasCalled).To(BeTrue())
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
				Expect(conn.BeginWasCalled).To(BeFalse())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeFalse())
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
				Expect(conn.BeginWasCalled).To(BeFalse())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeFalse())
			})

			It("delegates transaction errors to the error writer", func() {
				conn.CommitError = "transaction error, oh no"
				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(models.NewTransactionCommitError("transaction error, oh no")))
				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeTrue())
				Expect(conn.RollbackWasCalled).To(BeFalse())
			})
		})
	})
})
