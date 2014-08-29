package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/params"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Notify", func() {
    var handler handlers.Notify
    var fakeFinder *FakeFinder
    var fakeCourier *FakeCourier
    var fakeRegistrar *FakeRegistrar
    var request *http.Request
    var token string
    var client models.Client
    var kind models.Kind
    var fakeDBConn *FakeDBConn

    BeforeEach(func() {
        client = models.Client{
            ID:          "mister-client",
            Description: "Health Monitor",
        }
        kind = models.Kind{
            ID:          "test_email",
            Description: "Instance Down",
            ClientID:    "mister-client",
        }
        fakeFinder = NewFakeFinder()
        fakeFinder.Clients["mister-client"] = client
        fakeFinder.Kinds["test_email|mister-client"] = kind
        fakeCourier = NewFakeCourier()

        fakeRegistrar = NewFakeRegistrar()

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

        tokenHeader := map[string]interface{}{
            "alg": "FAST",
        }
        tokenClaims := map[string]interface{}{
            "client_id": "mister-client",
            "exp":       3404281214,
            "scope":     []string{"notifications.write"},
        }
        token = BuildToken(tokenHeader, tokenClaims)

        request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
        if err != nil {
            panic(err)
        }
        request.Header.Set("Authorization", "Bearer "+token)

        handler = handlers.NewNotify(fakeCourier, fakeFinder, fakeRegistrar)
    })

    It("delegates to the courier", func() {
        _, err := handler.Execute(fakeDBConn, request, postal.SpaceGUID("space-001"))
        if err != nil {
            panic(err)
        }

        Expect(fakeCourier.DispatchArguments).To(Equal([]interface{}{
            "mister-client",
            postal.SpaceGUID("space-001"),
            postal.Options{
                ReplyTo:           "me@example.com",
                Subject:           "Your instance is down",
                KindDescription:   "Instance Down",
                SourceDescription: "Health Monitor",
                Text:              "This is the plain text body of the email",
                HTML:              "<p>This is the HTML Body of the email</p>",
                KindID:            "test_email",
            },
        }))
    })

    It("registers the client and kind", func() {
        _, err := handler.Execute(fakeDBConn, request, postal.SpaceGUID("space-001"))
        if err != nil {
            panic(err)
        }

        Expect(fakeRegistrar.RegisterArguments).To(Equal([]interface{}{
            fakeDBConn,
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
                request.Header.Set("Authorization", "Bearer "+token)

                _, err = handler.Execute(fakeDBConn, request, postal.SpaceGUID("space-001"))

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
                request.Header.Set("Authorization", "Bearer "+token)

                _, err = handler.Execute(fakeDBConn, request, postal.SpaceGUID("space-001"))

                Expect(err).To(Equal(params.ParseError{}))
            })
        })

        Context("when the courier returns errors", func() {
            It("returns the error", func() {
                fakeCourier.Error = errors.New("BOOM!")

                _, err := handler.Execute(fakeDBConn, request, postal.UserGUID("user-123"))

                Expect(err).To(Equal(errors.New("BOOM!")))
            })
        })

        Context("when the finder return errors", func() {
            It("returns the error", func() {
                fakeFinder.ClientAndKindError = errors.New("BOOM!")

                _, err := handler.Execute(fakeDBConn, request, postal.UserGUID("user-123"))

                Expect(err).To(Equal(errors.New("BOOM!")))
            })
        })

        Context("when the registrar returns errors", func() {
            It("returns the error", func() {
                fakeRegistrar.RegisterError = errors.New("BOOM!")

                _, err := handler.Execute(fakeDBConn, request, postal.UserGUID("user-123"))

                Expect(err).To(Equal(errors.New("BOOM!")))
            })
        })
    })
})
