package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Receipts Repo", func() {
	var repo models.ReceiptsRepo
	var conn *db.Connection

	BeforeEach(func() {
		repo = models.NewReceiptsRepo()

		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)

		conn = database.Connection().(*db.Connection)
	})

	Describe("CreateReceipts", func() {
		var (
			firstUserGUID  string
			secondUserGUID string
			userGUIDs      []string
			clientID       string
			kindID         string
		)

		BeforeEach(func() {
			firstUserGUID = "user-123"
			secondUserGUID = "user-456"
			userGUIDs = []string{firstUserGUID, secondUserGUID}
			clientID = "client-abc"
			kindID = "be-kind"
		})

		It("creates or updates a receipt for each user", func() {
			err := repo.CreateReceipts(conn, userGUIDs, clientID, kindID)
			if err != nil {
				panic(err)
			}

			firstReceipt, err := findReceipt(conn, firstUserGUID, clientID, kindID)
			if err != nil {
				panic(err)
			}

			secondReceipt, err := findReceipt(conn, secondUserGUID, clientID, kindID)
			if err != nil {
				panic(err)
			}

			Expect(firstReceipt.UserGUID).To(Equal(firstUserGUID))
			Expect(firstReceipt.ClientID).To(Equal(clientID))
			Expect(firstReceipt.KindID).To(Equal(kindID))
			Expect(firstReceipt.Count).To(Equal(1))

			Expect(secondReceipt.UserGUID).To(Equal(secondUserGUID))
			Expect(secondReceipt.ClientID).To(Equal(clientID))
			Expect(secondReceipt.KindID).To(Equal(kindID))
			Expect(secondReceipt.Count).To(Equal(1))
		})

		It("updates a receipt's count for a user for a given clientID and kindID", func() {
			receipt := models.Receipt{
				UserGUID: firstUserGUID,
				ClientID: clientID,
				KindID:   kindID,
			}

			_, err := createReceipt(conn, receipt)
			if err != nil {
				panic(err)
			}

			rowCount, err := conn.SelectInt("SELECT COUNT(*) FROM `receipts`")
			if err != nil {
				panic(err)
			}

			Expect(int(rowCount)).To(Equal(1))

			err = repo.CreateReceipts(conn, []string{firstUserGUID}, clientID, kindID)
			if err != nil {
				panic(err)
			}

			rowCount, err = conn.SelectInt("SELECT COUNT(*) FROM `receipts`")
			if err != nil {
				panic(err)
			}

			Expect(int(rowCount)).To(Equal(1))

			firstReceipt, err := findReceipt(conn, firstUserGUID, clientID, kindID)
			if err != nil {
				panic(err)
			}

			Expect(firstReceipt.UserGUID).To(Equal(firstUserGUID))
			Expect(firstReceipt.ClientID).To(Equal(clientID))
			Expect(firstReceipt.KindID).To(Equal(kindID))
			Expect(firstReceipt.Count).To(Equal(2))
		})

		It("does not update count and adds a row when clientID and kindID are different", func() {
			_, err := createReceipt(conn, models.Receipt{
				UserGUID: firstUserGUID,
				ClientID: clientID,
				KindID:   kindID,
			})

			rowCount, err := conn.SelectInt("SELECT COUNT(*) FROM `receipts`")
			if err != nil {
				panic(err)
			}

			Expect(int(rowCount)).To(Equal(1))

			err = repo.CreateReceipts(conn, []string{firstUserGUID}, "weird-client", kindID)
			if err != nil {
				panic(err)
			}

			rowCount, err = conn.SelectInt("SELECT COUNT(*) FROM `receipts`")
			if err != nil {
				panic(err)
			}
			Expect(int(rowCount)).To(Equal(2))

			err = repo.CreateReceipts(conn, []string{firstUserGUID}, clientID, "a-new-kind")
			if err != nil {
				panic(err)
			}

			rowCount, err = conn.SelectInt("SELECT COUNT(*) FROM `receipts`")
			if err != nil {
				panic(err)
			}
			Expect(int(rowCount)).To(Equal(3))

			firstReceipt, err := findReceipt(conn, firstUserGUID, clientID, kindID)
			if err != nil {
				panic(err)
			}

			differentClientReceipt, err := findReceipt(conn, firstUserGUID, "weird-client", kindID)
			if err != nil {
				panic(err)
			}

			differentKindReceipt, err := findReceipt(conn, firstUserGUID, clientID, "a-new-kind")
			if err != nil {
				panic(err)
			}

			Expect(firstReceipt.UserGUID).To(Equal(firstUserGUID))
			Expect(differentKindReceipt.UserGUID).To(Equal(firstUserGUID))
			Expect(differentClientReceipt.UserGUID).To(Equal(firstUserGUID))
			Expect(firstReceipt.Primary).ToNot(Equal(differentClientReceipt.Primary))
			Expect(firstReceipt.Primary).ToNot(Equal(differentKindReceipt.Primary))
		})
	})
})
