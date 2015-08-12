package services_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationsFinder", func() {
	var (
		finder      services.NotificationsFinder
		clientsRepo *fakes.ClientsRepo
		kindsRepo   *fakes.KindsRepo
		database    *fakes.Database
	)

	BeforeEach(func() {
		clientsRepo = fakes.NewClientsRepo()
		kindsRepo = fakes.NewKindsRepo()
		database = fakes.NewDatabase()
		finder = services.NewNotificationsFinder(clientsRepo, kindsRepo)
	})

	Describe("ClientAndKind", func() {
		var (
			raptors models.Client
			breach  models.Kind
		)

		BeforeEach(func() {
			raptors = models.Client{
				ID:        "raptors",
				CreatedAt: time.Now(),
			}
			clientsRepo.Clients["raptors"] = raptors

			breach = models.Kind{
				ID:        "perimeter_breach",
				ClientID:  "raptors",
				CreatedAt: time.Now(),
			}
			kindsRepo.Kinds[breach.ID+breach.ClientID] = breach
		})

		It("retrieves clients and kinds from the database", func() {
			client, kind, err := finder.ClientAndKind(database, "raptors", "perimeter_breach")
			if err != nil {
				panic(err)
			}

			Expect(client).To(Equal(raptors))
			Expect(kind).To(Equal(breach))
		})

		Context("when the client cannot be found", func() {
			It("returns an empty models.Client", func() {
				client, _, err := finder.ClientAndKind(database, "bad-client-id", "perimeter_breach")
				Expect(client).To(Equal(models.Client{
					ID: "bad-client-id",
				}))
				Expect(err).To(BeNil())
			})
		})

		Context("when the kind cannot be found", func() {
			It("returns an empty models.Kind", func() {
				client, kind, err := finder.ClientAndKind(database, "raptors", "bad-kind-id")
				Expect(client).To(Equal(raptors))
				Expect(kind).To(Equal(models.Kind{
					ID:       "bad-kind-id",
					ClientID: "raptors",
				}))
				Expect(err).To(BeNil())
			})
		})

		Context("when the repo returns an error other than RecordNotFoundError", func() {
			It("returns the error", func() {
				clientsRepo.FindCall.Error = errors.New("BOOM!")
				_, _, err := finder.ClientAndKind(database, "raptors", "perimeter_breach")
				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})

		Context("when the kinds repo returns an error other than RecordNotFoundError", func() {
			It("returns the error", func() {
				kindsRepo.FindError = errors.New("BOOM!")
				_, _, err := finder.ClientAndKind(database, "raptors", "perimeter_breach")
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

			clientsRepo.Clients = map[string]models.Client{
				"the-imitation-game": imitationGame,
				"big-hero-6":         bigHero6,
				"star-wars":          starWars,
			}

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

			kindsRepo.Kinds = map[string]models.Kind{
				"star-wars|multi-light-saber": multiSaber,
				"star-wars|millenium-falcon":  milleniumFalcon,
				"big-hero-6|robots":           robots,
			}
		})

		It("returns all clients with their associated notifications", func() {
			clients := []models.Client{}
			for _, client := range clientsRepo.Clients {
				clients = append(clients, client)
			}

			clientsRepo.AllClients = clients
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
