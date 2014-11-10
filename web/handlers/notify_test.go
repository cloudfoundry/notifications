package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/postal/strategies"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/dgrijalva/jwt-go"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Notify", func() {
    Describe("Execute", func() {
        Context("When Emailing a user or a group", func() {
            var handler handlers.Notify
            var finder *fakes.Finder
            var registrar *fakes.Registrar
            var request *http.Request
            var rawToken string
            var client models.Client
            var kind models.Kind
            var conn *fakes.DBConn
            var strategy *fakes.MailStrategy
            var context stack.Context
            var tokenHeader map[string]interface{}
            var tokenClaims map[string]interface{}

            BeforeEach(func() {
                client = models.Client{
                    ID:          "mister-client",
                    Description: "Health Monitor",
                }
                kind = models.Kind{
                    ID:          "test_email",
                    Description: "Instance Down",
                    ClientID:    "mister-client",
                    Critical:    true,
                }
                finder = fakes.NewFinder()
                finder.Clients["mister-client"] = client
                finder.Kinds["test_email|mister-client"] = kind

                registrar = fakes.NewRegistrar()

                body, err := json.Marshal(map[string]string{
                    "kind_id":  "test_email",
                    "text":     "This is the plain text body of the email",
                    "html":     "<p>This is the HTML Body of the email</p>",
                    "subject":  "Your instance is down",
                    "reply_to": "me@example.com",
                })
                if err != nil {
                    panic(err)
                }

                tokenHeader = map[string]interface{}{
                    "alg": "FAST",
                }
                tokenClaims = map[string]interface{}{
                    "client_id": "mister-client",
                    "exp":       int64(3404281214),
                    "scope":     []string{"notifications.write", "critical_notifications.write"},
                }
                rawToken = fakes.BuildToken(tokenHeader, tokenClaims)

                request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }
                request.Header.Set("Authorization", "Bearer "+rawToken)

                token, err := jwt.Parse(rawToken, func(*jwt.Token) ([]byte, error) {
                    return []byte(config.UAAPublicKey), nil
                })

                context = stack.NewContext()
                context.Set("token", token)

                conn = fakes.NewDBConn()

                handler = handlers.NewNotify(finder, registrar)
                strategy = fakes.NewMailStrategy()
            })

            Describe("Responses", func() {
                BeforeEach(func() {
                    strategy.Responses = []strategies.Response{
                        {
                            Status:         "delivered",
                            Recipient:      "user-123",
                            NotificationID: "123-456",
                        },
                    }
                })

                It("trim is called on the strategy", func() {
                    _, err := handler.Execute(conn, request, context, postal.UAAGUID(""), strategy)
                    if err != nil {
                        panic(err)
                    }

                    Expect(strategy.TrimCalled).To(Equal(true))
                })
            })

            It("delegates to the mailStrategy", func() {
                _, err := handler.Execute(conn, request, context, postal.UAAGUID("space-001"), strategy)
                if err != nil {
                    panic(err)
                }

                Expect(strategy.DispatchArguments).To(Equal([]interface{}{
                    "mister-client",
                    "space-001",
                    postal.Options{
                        ReplyTo:           "me@example.com",
                        Subject:           "Your instance is down",
                        KindDescription:   "Instance Down",
                        SourceDescription: "Health Monitor",
                        Text:              "This is the plain text body of the email",
                        HTML:              postal.HTML{BodyAttributes: "", BodyContent: "<p>This is the HTML Body of the email</p>"},
                        KindID:            "test_email",
                    },
                }))
            })

            It("registers the client and kind", func() {
                _, err := handler.Execute(conn, request, context, postal.UAAGUID("space-001"), strategy)
                if err != nil {
                    panic(err)
                }

                Expect(registrar.RegisterArguments).To(Equal([]interface{}{
                    conn,
                    client,
                    []models.Kind{kind},
                }))
            })

            Context("failure cases", func() {
                Context("when validating params", func() {
                    It("returns a error response when params are missing", func() {
                        body, err := json.Marshal(map[string]string{
                            "subject": "Your instance is down",
                        })
                        if err != nil {
                            panic(err)
                        }
                        request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
                        if err != nil {
                            panic(err)
                        }
                        request.Header.Set("Authorization", "Bearer "+rawToken)

                        _, err = handler.Execute(conn, request, context, postal.UAAGUID("space-001"), strategy)

                        Expect(err).ToNot(BeNil())
                        validationErr := err.(params.ValidationError)
                        Expect(validationErr.Errors()).To(ContainElement(`"kind_id" is a required field`))
                        Expect(validationErr.Errors()).To(ContainElement(`"text" or "html" fields must be supplied`))
                    })

                    It("returns a error response when params cannot be parsed", func() {
                        request, err := http.NewRequest("POST", "/spaces/space-001", strings.NewReader("this is not JSON"))
                        if err != nil {
                            panic(err)
                        }
                        request.Header.Set("Authorization", "Bearer "+rawToken)

                        _, err = handler.Execute(conn, request, context, postal.UAAGUID("space-001"), strategy)

                        Expect(err).To(Equal(params.ParseError{}))
                    })
                })

                Context("when the strategy dispatch method returns errors", func() {
                    It("returns the error", func() {
                        strategy.Error = errors.New("BOOM!")

                        _, err := handler.Execute(conn, request, context, postal.UAAGUID("user-123"), strategy)

                        Expect(err).To(Equal(errors.New("BOOM!")))
                    })
                })

                Context("when the finder return errors", func() {
                    It("returns the error", func() {
                        finder.ClientAndKindError = errors.New("BOOM!")

                        _, err := handler.Execute(conn, request, context, postal.UAAGUID("user-123"), strategy)

                        Expect(err).To(Equal(errors.New("BOOM!")))
                    })
                })

                Context("when the registrar returns errors", func() {
                    It("returns the error", func() {
                        registrar.RegisterError = errors.New("BOOM!")

                        _, err := handler.Execute(conn, request, context, postal.UAAGUID("user-123"), strategy)

                        Expect(err).To(Equal(errors.New("BOOM!")))
                    })
                })

                Context("when trying to send a critical notification without the correct scope", func() {
                    It("returns an error", func() {
                        tokenClaims["scope"] = []interface{}{"notifications.write"}
                        rawToken = fakes.BuildToken(tokenHeader, tokenClaims)
                        token, err := jwt.Parse(rawToken, func(*jwt.Token) ([]byte, error) {
                            return []byte(config.UAAPublicKey), nil
                        })

                        context.Set("token", token)

                        _, err = handler.Execute(conn, request, context, postal.UAAGUID("user-123"), strategy)

                        Expect(err).To(BeAssignableToTypeOf(postal.NewCriticalNotificationError("test_email")))
                    })
                })
            })
        })
    })
})
