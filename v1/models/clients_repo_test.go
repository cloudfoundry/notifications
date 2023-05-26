package models_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientsRepo", func() {
	var (
		repo models.ClientsRepo
		conn db.ConnectionInterface
	)

	BeforeEach(func() {
		repo = models.NewClientsRepo()
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)
		conn = database.Connection()
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

			firstClient, err := repo.Upsert(conn, client1)
			if err != nil {
				panic(err)
			}

			secondClient, err := repo.Upsert(conn, client2)
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
		Context("when the template id is meant to be updated", func() {
			It("updates the record in the database", func() {
				client := models.Client{
					ID:         "my-client",
					TemplateID: "my-template",
				}

				client, err := repo.Upsert(conn, client)
				if err != nil {
					panic(err)
				}

				client.ID = "my-client"
				client.Description = "My Client"
				client.TemplateID = "new-template"

				client, err = repo.Update(conn, client)
				Expect(err).NotTo(HaveOccurred())

				client, err = repo.Find(conn, "my-client")
				if err != nil {
					panic(err)
				}

				Expect(client.ID).To(Equal("my-client"))
				Expect(client.Description).To(Equal("My Client"))
				Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
				Expect(client.TemplateID).To(Equal("new-template"))
			})

			It("returns a record not found error when the record does not exist", func() {
				client := models.Client{
					ID:         "my-client",
					TemplateID: "my-template",
				}

				_, err := repo.Update(conn, client)
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Client with ID \"my-client\" could not be found")}))
			})
		})

		Context("when the template id is not meant to be updated", func() {
			It("Uses the existing templateID when the field is empty", func() {
				client := models.Client{
					ID:         "my-client",
					TemplateID: "my-template",
				}

				client, err := repo.Upsert(conn, client)
				if err != nil {
					panic(err)
				}

				client.TemplateID = models.DoNotSetTemplateID
				client.Description = "My Client"

				client, err = repo.Update(conn, client)
				Expect(err).NotTo(HaveOccurred())

				client, err = repo.Find(conn, "my-client")
				if err != nil {
					panic(err)
				}

				Expect(client.ID).To(Equal("my-client"))
				Expect(client.Description).To(Equal("My Client"))
				Expect(client.TemplateID).To(Equal("my-template"))
				Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			})

			It("returns a record not found error when the record does not exist", func() {
				client := models.Client{
					ID: "my-client",
				}

				_, err := repo.Update(conn, client)
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Client with ID \"my-client\" could not be found")}))
			})
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

				client, err := repo.Upsert(conn, client)
				Expect(err).NotTo(HaveOccurred())

				client = models.Client{
					ID:          "my-client",
					Description: "My Client",
				}

				client, err = repo.Upsert(conn, client)
				Expect(err).NotTo(HaveOccurred())

				Expect(client.ID).To(Equal("my-client"))
				Expect(client.Description).To(Equal("My Client"))
				Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			})
		})

		Context("when the record comes into existence after the Find, but before we create it", func() {
			It("updates the record in the database", func() {
				client := models.Client{
					ID:          "my-client",
					Description: "My Client",
					TemplateID:  "some-template-id",
				}

				conn := mocks.NewConnection()
				conn.InsertCall.Returns.Error = errors.New("Duplicate entry")

				_, err := repo.Upsert(conn, client)
				Expect(err).NotTo(HaveOccurred())
				Expect(conn.UpdateCall.Receives.List).To(HaveLen(1))
				Expect(conn.UpdateCall.Receives.List[0].(*models.Client).ID).To(Equal("my-client"))
			})
		})
	})

	Describe("FindAllByTemplateID", func() {
		It("returns a list of clients with the given template ID", func() {
			client1, err := repo.Upsert(conn, models.Client{
				ID:         "i-have-a-template",
				TemplateID: "some-template-id",
			})
			if err != nil {
				panic(err)
			}
			_, err = repo.Upsert(conn, models.Client{
				ID: "i-dont-have-a-template",
			})
			if err != nil {
				panic(err)
			}

			returnedClients, err := repo.FindAllByTemplateID(conn, "some-template-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedClients).To(HaveLen(1))
			Expect(returnedClients).To(ContainElement(client1))
		})
	})
})
