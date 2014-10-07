package services_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Registrar", func() {
    var registrar services.Registrar
    var fakeClientsRepo *fakes.FakeClientsRepo
    var fakeKindsRepo *fakes.FakeKindsRepo
    var fakeDBConn *fakes.FakeDBConn
    var kinds []models.Kind

    BeforeEach(func() {
        fakeClientsRepo = fakes.NewFakeClientsRepo()
        fakeKindsRepo = fakes.NewFakeKindsRepo()
        registrar = services.NewRegistrar(fakeClientsRepo, fakeKindsRepo)
        fakeDBConn = &fakes.FakeDBConn{}
    })

    Describe("Register", func() {
        It("stores the client and kind records in the database", func() {
            client := models.Client{
                ID:          "raptors",
                Description: "perimeter breech",
            }

            hungry := models.Kind{
                ID:          "hungry",
                Description: "these raptors are hungry",
                Critical:    true,
                ClientID:    "raptors",
            }

            sleepy := models.Kind{
                ID:          "sleepy",
                Description: "these raptors are zzzzzzzz",
                Critical:    false,
                ClientID:    "raptors",
            }

            kinds = []models.Kind{hungry, sleepy}

            err := registrar.Register(fakeDBConn, client, kinds)
            if err != nil {
                panic(err)
            }

            Expect(len(fakeClientsRepo.Clients)).To(Equal(1))
            Expect(fakeClientsRepo.Clients["raptors"]).To(Equal(client))

            Expect(len(fakeKindsRepo.Kinds)).To(Equal(2))
            kind, err := fakeKindsRepo.Find(fakeDBConn, "hungry", "raptors")
            if err != nil {
                panic(err)
            }
            Expect(kind).To(Equal(hungry))

            kind, err = fakeKindsRepo.Find(fakeDBConn, "sleepy", "raptors")
            if err != nil {
                panic(err)
            }
            Expect(kind).To(Equal(sleepy))
        })

        It("idempotently updates the client and kinds", func() {
            _, err := fakeClientsRepo.Create(fakeDBConn, models.Client{
                ID:          "raptors",
                Description: "perimeter breech",
            })
            if err != nil {
                panic(err)
            }

            _, err = fakeKindsRepo.Create(fakeDBConn, models.Kind{
                ID:          "hungry",
                Description: "these raptors are hungry",
                Critical:    true,
                ClientID:    "raptors",
            })
            if err != nil {
                panic(err)
            }

            _, err = fakeKindsRepo.Create(fakeDBConn, models.Kind{
                ID:          "sleepy",
                Description: "these raptors are zzzzzzzz",
                Critical:    false,
                ClientID:    "raptors",
            })
            if err != nil {
                panic(err)
            }

            client := models.Client{
                ID:          "raptors",
                Description: "perimeter breech new descrition",
            }

            hungry := models.Kind{
                ID:          "hungry",
                Description: "these raptors are hungry new descrition",
                Critical:    true,
                ClientID:    "raptors",
            }

            sleepy := models.Kind{
                ID:          "sleepy",
                Description: "these raptors are zzzzzzzz new descrition",
                Critical:    false,
                ClientID:    "raptors",
            }

            kinds := []models.Kind{hungry, sleepy}

            err = registrar.Register(fakeDBConn, client, kinds)
            if err != nil {
                panic(err)
            }

            Expect(len(fakeClientsRepo.Clients)).To(Equal(1))
            Expect(fakeClientsRepo.Clients["raptors"]).To(Equal(client))

            Expect(len(fakeKindsRepo.Kinds)).To(Equal(2))
            kind, err := fakeKindsRepo.Find(fakeDBConn, "hungry", "raptors")
            if err != nil {
                panic(err)
            }
            Expect(kind).To(Equal(hungry))

            kind, err = fakeKindsRepo.Find(fakeDBConn, "sleepy", "raptors")
            if err != nil {
                panic(err)
            }
            Expect(kind).To(Equal(sleepy))

        })

        Context("error cases", func() {
            It("returns the errors from the clients repo", func() {
                fakeClientsRepo.UpsertError = errors.New("BOOM!")

                err := registrar.Register(fakeDBConn, models.Client{}, []models.Kind{})

                Expect(err).To(Equal(errors.New("BOOM!")))
            })

            It("returns the errors from the kinds repo", func() {
                fakeKindsRepo.UpsertError = errors.New("BOOM!")

                err := registrar.Register(fakeDBConn, models.Client{}, []models.Kind{models.Kind{}})

                Expect(err).To(Equal(errors.New("BOOM!")))
            })
        })
    })

    Describe("Prune", func() {
        It("Removes kinds from the database that are not passed in", func() {
            client, err := fakeClientsRepo.Create(fakeDBConn, models.Client{
                ID:          "raptors",
                Description: "perimeter breech",
            })
            if err != nil {
                panic(err)
            }

            kind, err := fakeKindsRepo.Create(fakeDBConn, models.Kind{
                ID:          "hungry",
                Description: "these raptors are hungry",
                Critical:    true,
                ClientID:    "raptors",
            })
            if err != nil {
                panic(err)
            }

            _, err = fakeKindsRepo.Create(fakeDBConn, models.Kind{
                ID:          "sleepy",
                Description: "these raptors are zzzzzzzz",
                Critical:    false,
                ClientID:    "raptors",
            })
            if err != nil {
                panic(err)
            }

            err = registrar.Prune(fakeDBConn, client, []models.Kind{kind})
            if err != nil {
                panic(err)
            }

            Expect(fakeKindsRepo.TrimArguments).To(Equal([]interface{}{client.ID, []string{"hungry"}}))
        })
    })
})
