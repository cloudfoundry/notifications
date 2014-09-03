package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Connection", func() {
    var conn *models.Connection

    BeforeEach(func() {
        conn = &models.Connection{}
    })

    Describe("Transaction", func() {
        It("returns an uninitialized transaction", func() {
            transaction := conn.Transaction()
            Expect(transaction).To(BeAssignableToTypeOf(&models.Transaction{}))
        })
    })
})
