package models_test

import (
    "time"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/coopernurse/gorp"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UnsubscribesRepo", func() {
    var repo models.UnsubscribesRepoInterface
    var conn *gorp.DbMap

    BeforeEach(func() {
        TruncateTables()
        repo = models.NewUnsubscribesRepo()
        conn = models.Database().Connection
    })
    Describe("Create", func() {
        It("stores the unsubscribe record into the database", func() {

            unsubscribe := models.Unsubscribe{
                ClientID: "raptors",
                KindID:   "hungry-kind",
                UserID:   "correct-user",
            }

            unsubscribe, err := repo.Create(conn, unsubscribe)
            if err != nil {
                panic(err)
            }

            unsubscribe, err = repo.Find(conn, "raptors", "hungry-kind", "correct-user")
            if err != nil {
                panic(err)
            }

            Expect(unsubscribe.ClientID).To(Equal("raptors"))
            Expect(unsubscribe.KindID).To(Equal("hungry-kind"))
            Expect(unsubscribe.UserID).To(Equal("correct-user"))
            Expect(unsubscribe.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
        })

        It("returns an error duplicate record if the record is already in the database", func() {
            unsubscribe := models.Unsubscribe{
                ClientID: "raptors",
                KindID:   "hungry-kind",
                UserID:   "correct-user",
            }

            _, err := repo.Create(conn, unsubscribe)
            if err != nil {
                panic(err)
            }

            _, err = repo.Create(conn, unsubscribe)
            Expect(err).To(BeAssignableToTypeOf(models.ErrDuplicateRecord{}))

        })

    })
})
