package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registrar", func() {
	var (
		registrar   services.Registrar
		clientsRepo *mocks.ClientsRepository
		kindsRepo   *mocks.KindsRepo
		conn        *mocks.Connection
		kinds       []models.Kind
	)

	BeforeEach(func() {
		clientsRepo = mocks.NewClientsRepository()
		kindsRepo = mocks.NewKindsRepo()
		registrar = services.NewRegistrar(clientsRepo, kindsRepo)
		conn = mocks.NewConnection()
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

			err := registrar.Register(conn, client, kinds)
			Expect(err).NotTo(HaveOccurred())

			Expect(clientsRepo.UpsertCall.Receives.Connection).To(Equal(conn))
			Expect(clientsRepo.UpsertCall.Receives.Client).To(Equal(client))

			Expect(kindsRepo.UpsertCall.Receives.Kinds).To(Equal([]models.Kind{hungry, sleepy}))
		})

		Context("when kinds is an empty set", func() {
			It("does nothing", func() {
				err := registrar.Register(conn, models.Client{}, []models.Kind{{}})
				Expect(err).ToNot(HaveOccurred())
				Expect(kindsRepo.UpsertCall.Receives.Kinds).To(HaveLen(0))
			})
		})

		Context("error cases", func() {
			It("returns the errors from the clients repo", func() {
				clientsRepo.UpsertCall.Returns.Error = errors.New("BOOM!")

				err := registrar.Register(conn, models.Client{}, []models.Kind{})
				Expect(err).To(MatchError(errors.New("BOOM!")))
			})

			It("returns the errors from the kinds repo", func() {
				kindsRepo.UpsertCall.Returns.Error = errors.New("BOOM!")

				err := registrar.Register(conn, models.Client{}, []models.Kind{
					{ID: "something"},
				})
				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})

	Describe("Prune", func() {
		It("Removes kinds from the database that are not passed in", func() {
			client := models.Client{
				ID:          "raptors",
				Description: "perimeter breech",
			}

			kind := models.Kind{
				ID:          "hungry",
				Description: "these raptors are hungry",
				Critical:    true,
				ClientID:    "raptors",
			}

			err := registrar.Prune(conn, client, []models.Kind{kind})
			Expect(err).NotTo(HaveOccurred())

			Expect(kindsRepo.TrimCall.Receives.Connection).To(Equal(conn))
			Expect(kindsRepo.TrimCall.Receives.ClientID).To(Equal(client.ID))
			Expect(kindsRepo.TrimCall.Receives.KindIDs).To(Equal([]string{"hungry"}))
		})
	})
})
