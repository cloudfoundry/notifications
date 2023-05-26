package services_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationsFinder", func() {
	var (
		finder      services.NotificationsFinder
		clientsRepo *mocks.ClientsRepository
		kindsRepo   *mocks.KindsRepo
		database    *mocks.Database
	)

	BeforeEach(func() {
		clientsRepo = mocks.NewClientsRepository()
		kindsRepo = mocks.NewKindsRepo()
		database = mocks.NewDatabase()
		finder = services.NewNotificationsFinder(clientsRepo, kindsRepo)
	})

	Describe("ClientAndKind", func() {
		BeforeEach(func() {
			clientsRepo.FindCall.Returns.Client = models.Client{
				ID: "some-client-id",
			}

			kindsRepo.FindCall.Returns.Kinds = []models.Kind{
				{
					ID:       "perimeter_breach",
					ClientID: "some-client-id",
				},
			}
		})

		It("retrieves clients and kinds from the database", func() {
			client, kind, err := finder.ClientAndKind(database, "some-client-id", "perimeter_breach")
			Expect(err).NotTo(HaveOccurred())
			Expect(client).To(Equal(models.Client{
				ID: "some-client-id",
			}))
			Expect(kind).To(Equal(models.Kind{
				ID:       "perimeter_breach",
				ClientID: "some-client-id",
			}))
		})

		Context("when the client cannot be found", func() {
			It("returns an empty models.Client", func() {
				clientsRepo.FindCall.Returns.Error = models.NotFoundError{Err: errors.New("not found")}

				client, _, err := finder.ClientAndKind(database, "missing-client-id", "perimeter_breach")
				Expect(err).NotTo(HaveOccurred())
				Expect(client).To(Equal(models.Client{
					ID: "missing-client-id",
				}))
			})
		})

		Context("when the kind cannot be found", func() {
			It("returns an empty models.Kind", func() {
				kindsRepo.FindCall.Returns.Error = models.NotFoundError{Err: errors.New("not found")}

				client, kind, err := finder.ClientAndKind(database, "some-client-id", "bad-kind-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(client).To(Equal(models.Client{
					ID: "some-client-id",
				}))
				Expect(kind).To(Equal(models.Kind{
					ID:       "bad-kind-id",
					ClientID: "some-client-id",
				}))
			})
		})

		Context("when the repo returns an error other than RecordNotFoundError", func() {
			It("returns the error", func() {
				clientsRepo.FindCall.Returns.Error = errors.New("BOOM!")

				_, _, err := finder.ClientAndKind(database, "some-client-id", "perimeter_breach")
				Expect(err).To(MatchError(errors.New("BOOM!")))
			})
		})

		Context("when the kinds repo returns an error other than RecordNotFoundError", func() {
			It("returns the error", func() {
				kindsRepo.FindCall.Returns.Error = errors.New("BOOM!")

				_, _, err := finder.ClientAndKind(database, "some-client-id", "perimeter_breach")
				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})

	Describe("AllClientsAndNotifications", func() {
		var (
			starWars        models.Client
			bigHero6        models.Client
			imitationGame   models.Client
			multiSaber      models.Kind
			milleniumFalcon models.Kind
			robots          models.Kind
		)

		BeforeEach(func() {
			starWars = models.Client{
				ID:          "star-wars",
				Description: "The Force Awakens",
				CreatedAt:   time.Now(),
			}
			bigHero6 = models.Client{
				ID:          "big-hero-6",
				Description: "Marvel",
				CreatedAt:   time.Now(),
			}
			imitationGame = models.Client{
				ID:          "the-imitation-game",
				Description: "Alan Turing",
				CreatedAt:   time.Now(),
			}

			clientsRepo.FindAllCall.Returns.Clients = []models.Client{imitationGame, bigHero6, starWars}

			multiSaber = models.Kind{
				ID:          "multi-light-saber",
				ClientID:    "star-wars",
				Description: "LOL WUT?",
				Critical:    false,
				CreatedAt:   time.Now(),
			}
			milleniumFalcon = models.Kind{
				ID:          "millenium-falcon",
				ClientID:    "star-wars",
				Description: "Awesome!",
				Critical:    true,
				CreatedAt:   time.Now(),
			}

			robots = models.Kind{
				ID:          "robots",
				ClientID:    "big-hero-6",
				Description: "hero",
				Critical:    true,
				CreatedAt:   time.Now(),
			}

			kindsRepo.FindAllCall.Returns.Kinds = []models.Kind{multiSaber, milleniumFalcon, robots}
		})

		It("returns all clients with their associated notifications", func() {
			clients, notifications, err := finder.AllClientsAndNotifications(database)
			Expect(err).NotTo(HaveOccurred())
			Expect(clients).To(HaveLen(3))
			Expect(clients).To(ContainElement(starWars))
			Expect(clients).To(ContainElement(bigHero6))
			Expect(clients).To(ContainElement(imitationGame))

			Expect(notifications).To(HaveLen(3))
			Expect(notifications).To(ContainElement(multiSaber))
			Expect(notifications).To(ContainElement(milleniumFalcon))
			Expect(notifications).To(ContainElement(robots))
		})
	})
})
