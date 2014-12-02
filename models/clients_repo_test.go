package models_test

import (
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientsRepo", func() {
	var repo models.ClientsRepo
	var conn models.ConnectionInterface

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewClientsRepo()
		env := config.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		conn = models.NewDatabase(env.DatabaseURL, migrationsPath).Connection()
	})

	Describe("Create", func() {
		It("stores the client record into the database", func() {
			client := models.Client{
				ID:          "my-client",
				Description: "My Client",
			}

			client, err := repo.Create(conn, client)
			if err != nil {
				panic(err)
			}

			client, err = repo.Find(conn, "my-client")
			if err != nil {
				panic(err)
			}

			Expect(client.ID).To(Equal("my-client"))
			Expect(client.Description).To(Equal("My Client"))
			Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
		})
	})

	Describe("FindAll", func() {
		It("returns all the records in the database", func() {
			client1 := models.Client{
				ID:          "client1",
				Description: "client1-description",
			}

			client2 := models.Client{
				ID:          "client2",
				Description: "client2-description",
			}

			firstClient, err := repo.Create(conn, client1)
			if err != nil {
				panic(err)
			}

			secondClient, err := repo.Create(conn, client2)
			if err != nil {
				panic(err)
			}

			clients, err := repo.FindAll(conn)
			if err != nil {
				panic(err)
			}

			Expect(clients).To(Equal([]models.Client{firstClient, secondClient}))
		})
	})

	Describe("Update", func() {
		It("updates the record in the database", func() {
			client := models.Client{
				ID: "my-client",
			}

			client, err := repo.Create(conn, client)
			if err != nil {
				panic(err)
			}

			client.ID = "my-client"
			client.Description = "My Client"

			client, err = repo.Update(conn, client)
			if err != nil {
				panic(err)
			}

			client, err = repo.Find(conn, "my-client")
			if err != nil {
				panic(err)
			}

			Expect(client.ID).To(Equal("my-client"))
			Expect(client.Description).To(Equal("My Client"))
			Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
		})
	})

	Describe("Upsert", func() {
		Context("when the record is new", func() {
			It("inserts the record in the database", func() {
				client := models.Client{
					ID:          "my-client",
					Description: "My Client",
				}

				client, err := repo.Upsert(conn, client)
				if err != nil {
					panic(err)
				}

				Expect(client.ID).To(Equal("my-client"))
				Expect(client.Description).To(Equal("My Client"))
				Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			})
		})

		Context("when the record exists", func() {
			It("updates the record in the database", func() {
				client := models.Client{
					ID: "my-client",
				}

				client, err := repo.Create(conn, client)
				if err != nil {
					panic(err)
				}

				client = models.Client{
					ID:          "my-client",
					Description: "My Client",
				}

				client, err = repo.Upsert(conn, client)
				if err != nil {
					panic(err)
				}

				Expect(client.ID).To(Equal("my-client"))
				Expect(client.Description).To(Equal("My Client"))
				Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			})
		})
	})
})
