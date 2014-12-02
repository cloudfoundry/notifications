package services_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationsFinder", func() {
	var finder services.NotificationsFinder
	var clientsRepo *fakes.ClientsRepo
	var kindsRepo *fakes.KindsRepo

	BeforeEach(func() {
		clientsRepo = fakes.NewClientsRepo()
		kindsRepo = fakes.NewKindsRepo()
		finder = services.NewNotificationsFinder(clientsRepo, kindsRepo, fakes.NewDatabase())
	})

	Describe("ClientAndKind", func() {
		var raptors models.Client
		var breach models.Kind

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

	Describe("AllClientNotifications", func() {
		BeforeEach(func() {
			clientsRepo.Clients["star-wars"] = models.Client{
				ID:          "star-wars",
				Description: "The Force Awakens",
				CreatedAt:   time.Now(),
			}
			clientsRepo.Clients["big-hero-6"] = models.Client{
				ID:          "big-hero-6",
				Description: "Marvel",
				CreatedAt:   time.Now(),
			}
			clientsRepo.Clients["the-imitation-game"] = models.Client{
				ID:          "the-imitation-game",
				Description: "Alan Turing",
				CreatedAt:   time.Now(),
			}

			kindsRepo.Kinds["star-wars|multi-light-saber"] = models.Kind{
				ID:          "multi-light-saber",
				ClientID:    "star-wars",
				Description: "LOL WUT?",
				Critical:    false,
				CreatedAt:   time.Now(),
			}
			kindsRepo.Kinds["star-wars|millenium-falcon"] = models.Kind{
				ID:          "millenium-falcon",
				ClientID:    "star-wars",
				Description: "Awesome!",
				Critical:    true,
				CreatedAt:   time.Now(),
			}

			kindsRepo.Kinds["big-hero-6|robots"] = models.Kind{
				ID:          "robots",
				ClientID:    "big-hero-6",
				Description: "hero",
				Critical:    true,
				CreatedAt:   time.Now(),
			}
		})

		It("returns all clients with their associated notifications", func() {
			clients := []models.Client{}
			for _, client := range clientsRepo.Clients {
				clients = append(clients, client)
			}

			clientsRepo.AllClients = clients
			clientsWithNotifications, err := finder.AllClientNotifications()
			Expect(err).NotTo(HaveOccurred())
			Expect(clientsWithNotifications).To(Equal(map[string]services.ClientWithNotifications{
				"star-wars": {
					Name: "The Force Awakens",
					Notifications: map[string]services.Notification{
						"multi-light-saber": {
							Description: "LOL WUT?",
							Critical:    false,
						},
						"millenium-falcon": {
							Description: "Awesome!",
							Critical:    true,
						},
					},
				},
				"big-hero-6": {
					Name: "Marvel",
					Notifications: map[string]services.Notification{
						"robots": {
							Description: "hero",
							Critical:    true,
						},
					},
				},
				"the-imitation-game": {
					Name:          "Alan Turing",
					Notifications: map[string]services.Notification{},
				},
			}))
		})
	})
})
