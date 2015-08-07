package preferences_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateUserPreferencesHandler", func() {
	Describe("Execute", func() {
		var (
			handler     preferences.UpdateUserPreferencesHandler
			writer      *httptest.ResponseRecorder
			request     *http.Request
			connection  *fakes.Connection
			context     stack.Context
			updater     *fakes.PreferenceUpdater
			userGUID    string
			errorWriter *fakes.ErrorWriter
		)

		BeforeEach(func() {
			connection = fakes.NewConnection()
			database := fakes.NewDatabase()
			database.Conn = connection

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

			userGUID = "the-correct-user"
			request, err = http.NewRequest("PATCH", "domain/user_preferences/"+userGUID, bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			tokenHeader := map[string]interface{}{
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"client_id": "mister-client",
				"exp":       int64(3404281214),
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

			updater = fakes.NewPreferenceUpdater()
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()

			handler = preferences.NewUpdateUserPreferencesHandler(updater, errorWriter)
		})

		It("Passes the correct arguments to PreferenceUpdater Execute", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(len(updater.ExecuteArguments)).To(Equal(4))

			Expect(reflect.ValueOf(updater.ExecuteArguments[0]).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
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
			Expect(updater.ExecuteArguments[3]).To(Equal(userGUID))
		})

		It("Returns a 204 status code when the Preference object does not error", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("Failure cases", func() {
			Context("when global_unsubscribe is not set", func() {
				It("returns an error when the clients key is missing", func() {
					userGUID = "the-correct-user"
					jsonBody := `{"raptor-client": {"containment-unit-breach": {"email": false}}}`

					request, err := http.NewRequest("PATCH", "domain/user_preferences/"+userGUID, bytes.NewBuffer([]byte(jsonBody)))
					Expect(err).NotTo(HaveOccurred())

					tokenHeader := map[string]interface{}{
						"alg": "FAST",
					}
					tokenClaims := map[string]interface{}{
						"client_id": "mister-client",
						"exp":       int64(3404281214),
					}
					rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
					request.Header.Set("Authorization", "Bearer "+rawToken)

					handler.ServeHTTP(writer, request, context)

					Expect(errorWriter.Error).ToNot(BeNil())
					Expect(errorWriter.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))
					Expect(connection.BeginWasCalled).To(BeFalse())
					Expect(connection.CommitWasCalled).To(BeFalse())
					Expect(connection.RollbackWasCalled).To(BeFalse())
				})
			})

			It("delegates MissingKindOrClientErrors as webutil.ValidationError to the ErrorWriter", func() {
				updater.ExecuteError = services.MissingKindOrClientError("BOOM!")

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.Error).To(Equal(webutil.ValidationError([]string{"BOOM!"})))

				Expect(connection.BeginWasCalled).To(BeTrue())
				Expect(connection.CommitWasCalled).To(BeFalse())
				Expect(connection.RollbackWasCalled).To(BeTrue())
			})

			It("delegates CriticalKindErrors as webutil.ValidationError to the ErrorWriter", func() {
				updater.ExecuteError = services.CriticalKindError("BOOM!")

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.Error).To(Equal(webutil.ValidationError([]string{"BOOM!"})))

				Expect(connection.BeginWasCalled).To(BeTrue())
				Expect(connection.CommitWasCalled).To(BeFalse())
				Expect(connection.RollbackWasCalled).To(BeTrue())
			})

			It("delegates other errors to the ErrorWriter", func() {
				updater.ExecuteError = errors.New("BOOM!")

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))

				Expect(connection.BeginWasCalled).To(BeTrue())
				Expect(connection.CommitWasCalled).To(BeFalse())
				Expect(connection.RollbackWasCalled).To(BeTrue())
			})

			It("delegates transaction errors to the error writer", func() {
				connection.CommitError = "transaction error!!!"

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.Error).To(BeAssignableToTypeOf(models.NewTransactionCommitError("transaction error!!!")))

				Expect(connection.BeginWasCalled).To(BeTrue())
				Expect(connection.CommitWasCalled).To(BeTrue())
				Expect(connection.RollbackWasCalled).To(BeFalse())
			})
		})
	})
})
