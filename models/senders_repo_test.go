package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendersRepo", func() {
	var (
		repo models.SendersRepository
		conn models.ConnectionInterface
	)

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewSendersRepository(fakes.NewIncrementingGUIDGenerator().Generate)
		db := models.NewDatabase(sqlDB, models.Config{})
		db.Setup()
		conn = db.Connection()
	})

	Describe("Insert", func() {
		It("inserts the record into the database", func() {
			sender := models.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			sender, err := repo.Insert(conn, sender)
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(models.Sender{
				ID:       "deadbeef-aabb-ccdd-eeff-001122334455",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
		})

		It("returns a duplicate record error when the name and client_id are taken", func() {
			sender := models.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			_, err := repo.Insert(conn, sender)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Insert(conn, sender)
			Expect(err).To(MatchError(models.DuplicateRecordError{}))
		})
	})

	Describe("GetByClientIDAndName", func() {
		It("fetches the sender given a client_id and name", func() {
			createdSender, err := repo.Insert(conn, models.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			sender, err := repo.GetByClientIDAndName(conn, "some-client-id", "some-sender")
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(createdSender))
		})

		It("fails to fetch the sender given a non-existent client_id and name", func() {
			_, err := repo.Insert(conn, models.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByClientIDAndName(conn, "some-other-client-id", "some-sender")
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			Expect(err.Error()).To(Equal(`Record Not Found: Sender with client_id "some-other-client-id" and name "some-sender" could not be found`))
		})
	})
})
