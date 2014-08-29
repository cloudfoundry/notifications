package services_test

import (
    "errors"
    "time"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Finder", func() {
    var finder services.NotificationFinder
    var clientsRepo *FakeClientsRepo
    var kindsRepo *FakeKindsRepo
    var raptors models.Client
    var breach models.Kind

    Describe("ClientAndKind", func() {
        BeforeEach(func() {
            clientsRepo = NewFakeClientsRepo()
            raptors = models.Client{
                ID:        "raptors",
                CreatedAt: time.Now(),
            }
            clientsRepo.Clients["raptors"] = raptors

            kindsRepo = NewFakeKindsRepo()
            breach = models.Kind{
                ID:        "perimeter_breach",
                ClientID:  "raptors",
                CreatedAt: time.Now(),
            }
            kindsRepo.Kinds[breach.ID+breach.ClientID] = breach

            finder = services.NewNotificationFinder(clientsRepo, kindsRepo)
        })

        It("retrieves clients and kinds from the database", func() {
            client, kind, err := finder.ClientAndKind("raptors", "perimeter_breach")
            if err != nil {
                panic(err)
            }

            Expect(client).To(Equal(raptors))
            Expect(kind).To(Equal(breach))
        })

        Context("when the client cannot be found", func() {
            It("returns an empty models.Client", func() {
                client, _, err := finder.ClientAndKind("bad-client-id", "perimeter_breach")
                Expect(client).To(Equal(models.Client{
                    ID: "bad-client-id",
                }))
                Expect(err).To(BeNil())
            })
        })

        Context("when the kind cannot be found", func() {
            It("returns an empty models.Kind", func() {
                client, kind, err := finder.ClientAndKind("raptors", "bad-kind-id")
                Expect(client).To(Equal(raptors))
                Expect(kind).To(Equal(models.Kind{
                    ID:       "bad-kind-id",
                    ClientID: "raptors",
                }))
                Expect(err).To(BeNil())
            })
        })

        Context("when the repo returns an error other than ErrRecordNotFound", func() {
            It("returns the error", func() {
                clientsRepo.FindError = errors.New("BOOM!")
                _, _, err := finder.ClientAndKind("raptors", "perimeter_breach")
                Expect(err).To(Equal(errors.New("BOOM!")))
            })
        })

        Context("when the kinds repo returns an error other than ErrRecordNotFound", func() {
            It("returns the error", func() {
                kindsRepo.FindError = errors.New("BOOM!")
                _, _, err := finder.ClientAndKind("raptors", "perimeter_breach")
                Expect(err).To(Equal(errors.New("BOOM!")))
            })
        })
    })
})
