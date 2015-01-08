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
	var clientsRepo *fakes.ClientsRepo
	var kindsRepo *fakes.KindsRepo
	var conn *fakes.DBConn
	var kinds []models.Kind

	BeforeEach(func() {
		clientsRepo = fakes.NewClientsRepo()
		kindsRepo = fakes.NewKindsRepo()
		registrar = services.NewRegistrar(clientsRepo, kindsRepo)
		conn = fakes.NewDBConn()
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
			if err != nil {
				panic(err)
			}

			Expect(len(clientsRepo.Clients)).To(Equal(1))
			Expect(clientsRepo.Clients["raptors"]).To(Equal(client))

			Expect(len(kindsRepo.Kinds)).To(Equal(2))
			kind, err := kindsRepo.Find(conn, "hungry", "raptors")
			if err != nil {
				panic(err)
			}
			Expect(kind).To(Equal(hungry))

			kind, err = kindsRepo.Find(conn, "sleepy", "raptors")
			if err != nil {
				panic(err)
			}
			Expect(kind).To(Equal(sleepy))
		})

		It("idempotently updates the client and kinds", func() {
			_, err := clientsRepo.Create(conn, models.Client{
				ID:          "raptors",
				Description: "perimeter breech",
			})
			if err != nil {
				panic(err)
			}

			_, err = kindsRepo.Create(conn, models.Kind{
				ID:          "hungry",
				Description: "these raptors are hungry",
				Critical:    true,
				ClientID:    "raptors",
			})
			if err != nil {
				panic(err)
			}

			_, err = kindsRepo.Create(conn, models.Kind{
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

			err = registrar.Register(conn, client, kinds)
			if err != nil {
				panic(err)
			}

			Expect(len(clientsRepo.Clients)).To(Equal(1))
			Expect(clientsRepo.Clients["raptors"]).To(Equal(client))

			Expect(len(kindsRepo.Kinds)).To(Equal(2))
			kind, err := kindsRepo.Find(conn, "hungry", "raptors")
			if err != nil {
				panic(err)
			}
			Expect(kind).To(Equal(hungry))

			kind, err = kindsRepo.Find(conn, "sleepy", "raptors")
			if err != nil {
				panic(err)
			}
			Expect(kind).To(Equal(sleepy))
		})

		Context("when kind is an empty record", func() {
			It("does nothing", func() {
				err := registrar.Register(conn, models.Client{}, []models.Kind{models.Kind{}})
				Expect(err).ToNot(HaveOccurred())
				Expect(kindsRepo.Kinds).To(HaveLen(0))
			})
		})

		Context("error cases", func() {
			It("returns the errors from the clients repo", func() {
				clientsRepo.UpsertError = errors.New("BOOM!")

				err := registrar.Register(conn, models.Client{}, []models.Kind{})

				Expect(err).To(Equal(errors.New("BOOM!")))
			})

			It("returns the errors from the kinds repo", func() {
				kindsRepo.UpsertError = errors.New("BOOM!")

				err := registrar.Register(conn, models.Client{}, []models.Kind{
					{ID: "something"},
				})

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})

	Describe("Prune", func() {
		It("Removes kinds from the database that are not passed in", func() {
			client, err := clientsRepo.Create(conn, models.Client{
				ID:          "raptors",
				Description: "perimeter breech",
			})
			if err != nil {
				panic(err)
			}

			kind, err := kindsRepo.Create(conn, models.Kind{
				ID:          "hungry",
				Description: "these raptors are hungry",
				Critical:    true,
				ClientID:    "raptors",
			})
			if err != nil {
				panic(err)
			}

			_, err = kindsRepo.Create(conn, models.Kind{
				ID:          "sleepy",
				Description: "these raptors are zzzzzzzz",
				Critical:    false,
				ClientID:    "raptors",
			})
			if err != nil {
				panic(err)
			}

			err = registrar.Prune(conn, client, []models.Kind{kind})
			if err != nil {
				panic(err)
			}

			Expect(kindsRepo.TrimArguments).To(Equal([]interface{}{client.ID, []string{"hungry"}}))
		})
	})
})
