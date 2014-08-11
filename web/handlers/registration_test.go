package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

type FakeClientsRepo struct {
    Clients     map[string]models.Client
    UpsertError error
}

func NewFakeClientsRepo() *FakeClientsRepo {
    return &FakeClientsRepo{
        Clients: make(map[string]models.Client),
    }
}

func (fake *FakeClientsRepo) Create(client models.Client) (models.Client, error) {
    fake.Clients[client.ID] = client
    return client, nil
}

func (fake *FakeClientsRepo) Update(client models.Client) (models.Client, error) {
    fake.Clients[client.ID] = client
    return client, nil
}

func (fake *FakeClientsRepo) Upsert(client models.Client) (models.Client, error) {
    fake.Clients[client.ID] = client
    return client, fake.UpsertError
}

func (fake *FakeClientsRepo) Find(id string) (models.Client, error) {
    if client, ok := fake.Clients[id]; ok {
        return client, nil
    }
    return models.Client{}, models.ErrRecordNotFound{}
}

type FakeKindsRepo struct {
    Kinds       map[string]models.Kind
    UpsertError error
}

func NewFakeKindsRepo() *FakeKindsRepo {
    return &FakeKindsRepo{
        Kinds: make(map[string]models.Kind),
    }
}

func (fake *FakeKindsRepo) Create(kind models.Kind) (models.Kind, error) {
    fake.Kinds[kind.ID] = kind
    return kind, nil
}

func (fake *FakeKindsRepo) Update(kind models.Kind) (models.Kind, error) {
    fake.Kinds[kind.ID] = kind
    return kind, nil
}

func (fake *FakeKindsRepo) Upsert(kind models.Kind) (models.Kind, error) {
    fake.Kinds[kind.ID] = kind
    return kind, fake.UpsertError
}

func (fake *FakeKindsRepo) Find(id string) (models.Kind, error) {
    if kind, ok := fake.Kinds[id]; ok {
        return kind, nil
    }
    return models.Kind{}, models.ErrRecordNotFound{}
}

var _ = Describe("Registration", func() {
    var handler handlers.Registration
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var fakeClientsRepo *FakeClientsRepo
    var fakeKindsRepo *FakeKindsRepo
    var fakeErrorWriter *FakeErrorWriter

    BeforeEach(func() {
        fakeClientsRepo = NewFakeClientsRepo()
        fakeKindsRepo = NewFakeKindsRepo()
        fakeErrorWriter = &FakeErrorWriter{}
        handler = handlers.NewRegistration(fakeClientsRepo, fakeKindsRepo, fakeErrorWriter)
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
            "exp":       3404281214,
            "scope":     []string{"notifications.write"},
        }
        request.Header.Set("Authorization", "Bearer "+BuildToken(tokenHeader, tokenClaims))
    })

    It("stores the client and kind records in the database", func() {
        handler.ServeHTTP(writer, request)

        Expect(writer.Code).To(Equal(http.StatusOK))

        Expect(len(fakeClientsRepo.Clients)).To(Equal(1))
        client := fakeClientsRepo.Clients["raptors"]

        Expect(client.ID).To(Equal("raptors"))
        Expect(client.Description).To(Equal("Raptor Containment Unit"))

        Expect(len(fakeKindsRepo.Kinds)).To(Equal(2))
        Expect(fakeKindsRepo.Kinds["perimeter_breach"]).To(Equal(models.Kind{
            ID:          "perimeter_breach",
            Description: "Perimeter Breach",
            Critical:    true,
            ClientID:    "raptors",
        }))

        Expect(fakeKindsRepo.Kinds["feeding_time"]).To(Equal(models.Kind{
            ID:          "feeding_time",
            Description: "Feeding Time",
            Critical:    false,
            ClientID:    "raptors",
        }))
    })

    It("idempotently updates the client and kinds", func() {
        _, err := fakeClientsRepo.Create(models.Client{
            ID: "raptors",
        })
        if err != nil {
            panic(err)
        }

        _, err = fakeKindsRepo.Create(models.Kind{
            ID: "perimeter_breach",
        })
        if err != nil {
            panic(err)
        }

        _, err = fakeKindsRepo.Create(models.Kind{
            ID: "feeding_time",
        })
        if err != nil {
            panic(err)
        }

        handler.ServeHTTP(writer, request)

        Expect(writer.Code).To(Equal(http.StatusOK))

        Expect(len(fakeClientsRepo.Clients)).To(Equal(1))
        Expect(fakeClientsRepo.Clients["raptors"]).To(Equal(models.Client{
            ID:          "raptors",
            Description: "Raptor Containment Unit",
        }))

        Expect(len(fakeKindsRepo.Kinds)).To(Equal(2))
        Expect(fakeKindsRepo.Kinds["perimeter_breach"]).To(Equal(models.Kind{
            ID:          "perimeter_breach",
            Description: "Perimeter Breach",
            Critical:    true,
            ClientID:    "raptors",
        }))

        Expect(fakeKindsRepo.Kinds["feeding_time"]).To(Equal(models.Kind{
            ID:          "feeding_time",
            Description: "Feeding Time",
            Critical:    false,
            ClientID:    "raptors",
        }))
    })

    Context("failure cases", func() {
        It("delegates parsing errors to the ErrorWriter", func() {
            var err error
            request, err = http.NewRequest("PUT", "/registration", strings.NewReader("this is not valid JSON"))
            if err != nil {
                panic(err)
            }

            handler.ServeHTTP(writer, request)

            Expect(fakeErrorWriter.Error).To(BeAssignableToTypeOf(handlers.ParamsParseError{}))
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

            handler.ServeHTTP(writer, request)

            Expect(fakeErrorWriter.Error).To(BeAssignableToTypeOf(handlers.ParamsValidationError{}))
        })

        It("delegates client repo errors to the ErrorWriter", func() {
            fakeClientsRepo.UpsertError = models.ErrDuplicateRecord{}

            handler.ServeHTTP(writer, request)

            Expect(fakeErrorWriter.Error).To(BeAssignableToTypeOf(models.ErrDuplicateRecord{}))
        })

        It("delegates kind repo errors to the ErrorWriter", func() {
            fakeKindsRepo.UpsertError = models.ErrDuplicateRecord{}

            handler.ServeHTTP(writer, request)

            Expect(fakeErrorWriter.Error).To(BeAssignableToTypeOf(models.ErrDuplicateRecord{}))
        })
    })
})
