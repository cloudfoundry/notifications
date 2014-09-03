package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/params"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Registration", func() {
    var handler handlers.Registration
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var fakeErrorWriter *FakeErrorWriter
    var fakeConn *FakeDBConn
    var registrar *FakeRegistrar
    var client models.Client
    var kinds []models.Kind

    BeforeEach(func() {
        fakeConn = &FakeDBConn{}
        fakeErrorWriter = &FakeErrorWriter{}
        registrar = NewFakeRegistrar()
        handler = handlers.NewRegistration(registrar, fakeErrorWriter)
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
            "scope":     []string{"notifications.write"},
        }
        request.Header.Set("Authorization", "Bearer "+BuildToken(tokenHeader, tokenClaims))

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
            handler.Execute(writer, request, fakeConn)

            Expect(registrar.RegisterArguments).To(Equal([]interface{}{fakeConn, client, kinds}))

            Expect(fakeConn.BeginWasCalled).To(BeTrue())
            Expect(fakeConn.CommitWasCalled).To(BeTrue())
            Expect(fakeConn.RollbackWasCalled).To(BeFalse())
        })

        It("passes the correct arguments to Prune", func() {
            handler.Execute(writer, request, fakeConn)

            Expect(registrar.PruneArguments).To(Equal([]interface{}{fakeConn, client, kinds}))

            Expect(fakeConn.BeginWasCalled).To(BeTrue())
            Expect(fakeConn.CommitWasCalled).To(BeTrue())
            Expect(fakeConn.RollbackWasCalled).To(BeFalse())
        })

        It("does not trim kinds if they are not in the request", func() {
            requestBody, err := json.Marshal(map[string]interface{}{
                "source_description": "Raptor Containment Unit",
            })
            if err != nil {
                panic(err)
            }
            request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

            handler.Execute(writer, request, fakeConn)

            Expect(registrar.PruneArguments).To(BeNil())

            Expect(fakeConn.BeginWasCalled).To(BeTrue())
            Expect(fakeConn.CommitWasCalled).To(BeTrue())
            Expect(fakeConn.RollbackWasCalled).To(BeFalse())
        })

        Context("failure cases", func() {
            It("delegates parsing errors to the ErrorWriter", func() {
                var err error
                request, err = http.NewRequest("PUT", "/registration", strings.NewReader("this is not valid JSON"))
                if err != nil {
                    panic(err)
                }

                handler.Execute(writer, request, fakeConn)

                Expect(fakeErrorWriter.Error).To(BeAssignableToTypeOf(params.ParseError{}))

                Expect(fakeConn.BeginWasCalled).To(BeFalse())
                Expect(fakeConn.CommitWasCalled).To(BeFalse())
                Expect(fakeConn.RollbackWasCalled).To(BeFalse())
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

                handler.Execute(writer, request, fakeConn)

                Expect(fakeErrorWriter.Error).To(BeAssignableToTypeOf(params.ValidationError{}))

                Expect(fakeConn.BeginWasCalled).To(BeFalse())
                Expect(fakeConn.CommitWasCalled).To(BeFalse())
                Expect(fakeConn.RollbackWasCalled).To(BeFalse())
            })

            It("delegates registrar register errors to the ErrorWriter", func() {
                registrar.RegisterError = errors.New("BOOM!")

                handler.Execute(writer, request, fakeConn)

                Expect(fakeErrorWriter.Error).To(Equal(errors.New("BOOM!")))

                Expect(fakeConn.BeginWasCalled).To(BeTrue())
                Expect(fakeConn.CommitWasCalled).To(BeFalse())
                Expect(fakeConn.RollbackWasCalled).To(BeTrue())
            })

            It("delegates registrar prune errors to the ErrorWriter", func() {
                registrar.PruneError = errors.New("BOOM!")

                handler.Execute(writer, request, fakeConn)

                Expect(fakeErrorWriter.Error).To(Equal(errors.New("BOOM!")))

                Expect(fakeConn.BeginWasCalled).To(BeTrue())
                Expect(fakeConn.CommitWasCalled).To(BeFalse())
                Expect(fakeConn.RollbackWasCalled).To(BeTrue())
            })
        })
    })
})
