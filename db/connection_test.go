package db_test

import (
	"github.com/cloudfoundry-incubator/notifications/db"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection", func() {
	var conn *db.Connection

	BeforeEach(func() {
		conn = &db.Connection{}
	})

	Describe("Transaction", func() {
		It("returns an uninitialized transaction", func() {
			transaction := conn.Transaction()
			Expect(transaction).To(BeAssignableToTypeOf(&db.Transaction{}))
		})
	})
})
