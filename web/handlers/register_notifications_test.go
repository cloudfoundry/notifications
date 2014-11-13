package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegisterNotifications", func() {
	var handler handlers.RegisterNotifications
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var errorWriter *fakes.ErrorWriter
	var conn *fakes.DBConn
	var registrar *fakes.Registrar
	var client models.Client
	var kinds []models.Kind
	var context stack.Context

	BeforeEach(func() {
		conn = fakes.NewDBConn()
		errorWriter = fakes.NewErrorWriter()
		registrar = fakes.NewRegistrar()
		fakeDatabase := fakes.NewDatabase()
		handler = handlers.NewRegisterNotifications(registrar, errorWriter, fakeDatabase)
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
		if err != nil {
			panic(err)
		}

		request, err = http.NewRequest("PUT", "/registration", bytes.NewBuffer(requestBody))
		if err != nil {
			panic(err)
		}
		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "raptors",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write", "critical_notifications.write"},
		}
		rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
		request.Header.Set("Authorization", "Bearer "+rawToken)

		token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
			return []byte(config.UAAPublicKey), nil
		})
		context = stack.NewContext()
		context.Set("token", token)

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
	})

	Describe("Execute", func() {
		It("passes the correct arguments to Register", func() {
			handler.Execute(writer, request, conn, context)

			Expect(registrar.RegisterArguments).To(Equal([]interface{}{conn, client, kinds}))

			Expect(conn.BeginWasCalled).To(BeTrue())
			Expect(conn.CommitWasCalled).To(BeTrue())
			Expect(conn.RollbackWasCalled).To(BeFalse())
		})

		It("passes the correct arguments to Prune", func() {
			handler.Execute(writer, request, conn, context)

			Expect(registrar.PruneArguments).To(Equal([]interface{}{conn, client, kinds}))

			Expect(conn.BeginWasCalled).To(BeTrue())
			Expect(conn.CommitWasCalled).To(BeTrue())
			Expect(conn.RollbackWasCalled).To(BeFalse())
		})

		It("does not trim kinds if they are not in the request", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"source_description": "Raptor Containment Unit",
			})
			if err != nil {
				panic(err)
			}
			request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

			handler.Execute(writer, request, conn, context)

			Expect(registrar.PruneArguments).To(BeNil())

			Expect(conn.BeginWasCalled).To(BeTrue())
			Expect(conn.CommitWasCalled).To(BeTrue())
			Expect(conn.RollbackWasCalled).To(BeFalse())
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
				if err != nil {
					panic(err)
				}

				request, err = http.NewRequest("PUT", "/registration", bytes.NewBuffer(requestBody))
				if err != nil {
					panic(err)
				}
				tokenHeader := map[string]interface{}{
					"alg": "FAST",
				}
				tokenClaims := map[string]interface{}{
					"client_id": "raptors",
					"exp":       int64(3404281214),
					"scope":     []string{"notifications.write"},
				}
				rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
				request.Header.Set("Authorization", "Bearer "+rawToken)

				token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
					return []byte(config.UAAPublicKey), nil
				})
				if err != nil {
					panic(err)
				}

				context = stack.NewContext()
				context.Set("token", token)

				handler.Execute(writer, request, conn, context)
				Expect(errorWriter.Error).To(BeAssignableToTypeOf(postal.UAAScopesError("waaaaat")))

				Expect(conn.BeginWasCalled).To(BeFalse())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeFalse())
			})

			It("delegates parsing errors to the ErrorWriter", func() {
				var err error
				request, err = http.NewRequest("PUT", "/registration", strings.NewReader("this is not valid JSON"))
				if err != nil {
					panic(err)
				}

				handler.Execute(writer, request, conn, context)

				Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ParseError{}))

				Expect(conn.BeginWasCalled).To(BeFalse())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeFalse())
			})

			It("delegates validation errors to the ErrorWriter", func() {
				requestBody, err := json.Marshal(map[string]interface{}{})
				if err != nil {
					panic(err)
				}
				request, err = http.NewRequest("PUT", "/registration", bytes.NewBuffer(requestBody))
				if err != nil {
					panic(err)
				}

				handler.Execute(writer, request, conn, context)

				Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ValidationError{}))

				Expect(conn.BeginWasCalled).To(BeFalse())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeFalse())
			})

			It("delegates registrar register errors to the ErrorWriter", func() {
				registrar.RegisterError = errors.New("BOOM!")

				handler.Execute(writer, request, conn, context)

				Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeTrue())
			})

			It("delegates registrar prune errors to the ErrorWriter", func() {
				registrar.PruneError = errors.New("BOOM!")

				handler.Execute(writer, request, conn, context)

				Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))

				Expect(conn.BeginWasCalled).To(BeTrue())
				Expect(conn.CommitWasCalled).To(BeFalse())
				Expect(conn.RollbackWasCalled).To(BeTrue())
			})
		})
	})
})
