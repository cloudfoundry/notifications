package notifications_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegistrationHandler", func() {
	var (
		handler     notifications.RegistrationHandler
		writer      *httptest.ResponseRecorder
		request     *http.Request
		errorWriter *mocks.ErrorWriter
		conn        *mocks.Connection
		transaction *mocks.Transaction
		registrar   *mocks.Registrar
		client      models.Client
		kinds       []models.Kind
		context     stack.Context
	)

	BeforeEach(func() {
		transaction = mocks.NewTransaction()
		conn = mocks.NewConnection()
		conn.TransactionCall.Returns.Transaction = transaction
		database := mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		errorWriter = mocks.NewErrorWriter()
		registrar = mocks.NewRegistrar()
		writer = httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]interface{}{
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

		request, err = http.NewRequest("PUT", "/registration", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		tokenHeader := map[string]interface{}{
			"alg": "RS256",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "raptors",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write", "critical_notifications.write"},
		}
		rawToken := helpers.BuildToken(tokenHeader, tokenClaims)
		request.Header.Set("Authorization", "Bearer "+rawToken)

		token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
			return helpers.UAAPublicKeyRSA, nil
		})
		context = stack.NewContext()
		context.Set("token", token)
		context.Set("database", database)

		client = models.Client{
			ID:          "raptors",
			Description: "Raptor Containment Unit",
		}

		kinds = []models.Kind{
			{
				ID:          "perimeter_breach",
				Description: "Perimeter Breach",
				Critical:    true,
				ClientID:    client.ID,
			},
			{
				ID:          "feeding_time",
				Description: "Feeding Time",
				ClientID:    client.ID,
			},
		}

		handler = notifications.NewRegistrationHandler(registrar, errorWriter)
	})

	Describe("Execute", func() {
		It("passes the correct arguments to Register", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(registrar.RegisterCall.Receives.Connection).To(Equal(transaction))
			Expect(registrar.RegisterCall.Receives.Client).To(Equal(client))
			Expect(registrar.RegisterCall.Receives.Kinds).To(ConsistOf(kinds))

			Expect(transaction.BeginCall.WasCalled).To(BeTrue())
			Expect(transaction.CommitCall.WasCalled).To(BeTrue())
			Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
		})

		It("passes the correct arguments to Prune", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(registrar.PruneCall.Receives.Connection).To(Equal(transaction))
			Expect(registrar.PruneCall.Receives.Client).To(Equal(client))
			Expect(registrar.PruneCall.Receives.Kinds).To(ConsistOf(kinds))

			Expect(transaction.BeginCall.WasCalled).To(BeTrue())
			Expect(transaction.CommitCall.WasCalled).To(BeTrue())
			Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
		})

		It("does not trim kinds if they are not in the request", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"source_description": "Raptor Containment Unit",
			})
			Expect(err).NotTo(HaveOccurred())

			request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

			handler.ServeHTTP(writer, request, context)

			Expect(registrar.PruneCall.Called).To(BeFalse())

			Expect(transaction.BeginCall.WasCalled).To(BeTrue())
			Expect(transaction.CommitCall.WasCalled).To(BeTrue())
			Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
		})

		Context("failure cases", func() {
			It("rejects entire request and returns 404 error if notification is critical without scope", func() {
				requestBody, err := json.Marshal(map[string]interface{}{
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
							"critical":    true,
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("PUT", "/registration", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())

				tokenHeader := map[string]interface{}{
					"alg": "RS256",
				}
				tokenClaims := map[string]interface{}{
					"client_id": "raptors",
					"exp":       int64(3404281214),
					"scope":     []string{"notifications.write"},
				}
				rawToken := helpers.BuildToken(tokenHeader, tokenClaims)
				request.Header.Set("Authorization", "Bearer "+rawToken)

				token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
					return helpers.UAAPublicKeyRSA, nil
				})
				Expect(err).NotTo(HaveOccurred())

				context.Set("token", token)

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.UAAScopesError{Err: errors.New("UAA Scopes Error: Client does not have authority to register critical notifications.")}))

				Expect(transaction.BeginCall.WasCalled).To(BeFalse())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
			})

			It("delegates parsing errors to the ErrorWriter", func() {
				request, err := http.NewRequest("PUT", "/registration", strings.NewReader("this is not valid JSON"))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ParseError{}))

				Expect(transaction.BeginCall.WasCalled).To(BeFalse())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
			})

			It("delegates validation errors to the ErrorWriter", func() {
				requestBody, err := json.Marshal(map[string]interface{}{})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("PUT", "/registration", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))

				Expect(transaction.BeginCall.WasCalled).To(BeFalse())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())
			})

			It("delegates registrar register errors to the ErrorWriter", func() {
				registrar.RegisterCall.Returns.Error = errors.New("BOOM!")

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BOOM!")))

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
			})

			It("delegates registrar prune errors to the ErrorWriter", func() {
				registrar.PruneCall.Returns.Error = errors.New("BOOM!")

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BOOM!")))

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeFalse())
				Expect(transaction.RollbackCall.WasCalled).To(BeTrue())
			})

			It("delegates transaction errors to the ErrorWriter", func() {
				transaction.CommitCall.Returns.Error = errors.New("transaction commit error")
				handler.ServeHTTP(writer, request, context)

				Expect(transaction.BeginCall.WasCalled).To(BeTrue())
				Expect(transaction.CommitCall.WasCalled).To(BeTrue())
				Expect(transaction.RollbackCall.WasCalled).To(BeFalse())

				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("transaction commit error")))
			})
		})
	})
})
