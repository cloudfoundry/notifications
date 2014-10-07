package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UpdatePreferences", func() {
    Describe("Execute", func() {
        var handler handlers.UpdatePreferences
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var updater *fakes.FakePreferenceUpdater
        var errorWriter *fakes.FakeErrorWriter
        var fakeDBConn *fakes.FakeDBConn
        var context stack.Context

        BeforeEach(func() {
            fakeDBConn = &fakes.FakeDBConn{}
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
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "user_id": "correct-user",
                "exp":     int64(3404281214),
            }

            rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
            request.Header.Set("Authorization", "Bearer "+rawToken)

            token, err := jwt.Parse(rawToken, func(*jwt.Token) ([]byte, error) {
                return []byte(config.UAAPublicKey), nil
            })

            context = stack.NewContext()
            context.Set("token", token)

            errorWriter = fakes.NewFakeErrorWriter()
            updater = fakes.NewFakePreferenceUpdater()
            fakeDatabase := fakes.NewDatabase()
            handler = handlers.NewUpdatePreferences(updater, errorWriter, fakeDatabase)
            writer = httptest.NewRecorder()
        })

        It("Passes The Correct Arguments to PreferenceUpdater Execute", func() {
            handler.Execute(writer, request, fakeDBConn, context)
            Expect(len(updater.ExecuteArguments)).To(Equal(3))

            preferencesArguments := updater.ExecuteArguments[0]

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

            Expect(updater.ExecuteArguments[1]).To(BeTrue())
            Expect(updater.ExecuteArguments[2]).To(Equal("correct-user"))
        })

        It("Returns a 204 status code when the Preference object does not error", func() {
            handler.Execute(writer, request, fakeDBConn, context)

            Expect(writer.Code).To(Equal(http.StatusNoContent))
        })

        Context("Failure cases", func() {
            Context("preferenceUpdater.Execute errors", func() {

                It("delegates MissingKindOrClientErrors as params.ValidationError to the ErrorWriter", func() {
                    updater.ExecuteError = services.MissingKindOrClientError("BOOM!")

                    handler.Execute(writer, request, fakeDBConn, context)

                    Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{"BOOM!"})))

                    Expect(fakeDBConn.BeginWasCalled).To(BeTrue())
                    Expect(fakeDBConn.CommitWasCalled).To(BeFalse())
                    Expect(fakeDBConn.RollbackWasCalled).To(BeTrue())
                })

                It("delegates CriticalKindErrors as params.ValidationError to the ErrorWriter", func() {
                    updater.ExecuteError = services.CriticalKindError("BOOM!")

                    handler.Execute(writer, request, fakeDBConn, context)

                    Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{"BOOM!"})))

                    Expect(fakeDBConn.BeginWasCalled).To(BeTrue())
                    Expect(fakeDBConn.CommitWasCalled).To(BeFalse())
                    Expect(fakeDBConn.RollbackWasCalled).To(BeTrue())
                })

                It("delegates other errors to the ErrorWriter", func() {
                    updater.ExecuteError = errors.New("BOOM!")

                    handler.Execute(writer, request, fakeDBConn, context)

                    Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))

                    Expect(fakeDBConn.BeginWasCalled).To(BeTrue())
                    Expect(fakeDBConn.CommitWasCalled).To(BeFalse())
                    Expect(fakeDBConn.RollbackWasCalled).To(BeTrue())
                })
            })

            It("delegates parse errors to the ErrorWriter", func() {
                requestBody, err := json.Marshal([]string{})
                if err != nil {
                    panic(err)
                }

                request, err = http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer(requestBody))
                if err != nil {
                    panic(err)
                }

                handler.Execute(writer, request, fakeDBConn, context)

                Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ParseError{}))
                Expect(fakeDBConn.BeginWasCalled).To(BeFalse())
                Expect(fakeDBConn.CommitWasCalled).To(BeFalse())
                Expect(fakeDBConn.RollbackWasCalled).To(BeFalse())
            })

            It("delegates validation errors to the error writer", func() {
                requestBody, err := json.Marshal(map[string]map[string]map[string]map[string]interface{}{
                    "clients": {
                        "client-id": {
                            "kind-id": {},
                        },
                    },
                })
                if err != nil {
                    panic(err)
                }

                request, err = http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer(requestBody))
                if err != nil {
                    panic(err)
                }

                handler.Execute(writer, request, fakeDBConn, context)

                Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ValidationError{}))
                Expect(fakeDBConn.BeginWasCalled).To(BeFalse())
                Expect(fakeDBConn.CommitWasCalled).To(BeFalse())
                Expect(fakeDBConn.RollbackWasCalled).To(BeFalse())
            })
        })
    })

})
