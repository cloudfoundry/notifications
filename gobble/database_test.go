package gobble_test

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Column struct {
	Field string
	Type  string
}

var _ = Describe("Database", func() {
	It("has a jobs table", func() {
		database := gobble.NewDatabase(sqlDB)

		rows, err := database.Connection.Db.Query("SHOW TABLES")
		Expect(err).NotTo(HaveOccurred())

		defer rows.Close()
		tables := []string{}

		for rows.Next() {
			var table string
			err := rows.Scan(&table)
			Expect(err).NotTo(HaveOccurred())

			tables = append(tables, table)
		}

		Expect(tables).To(ContainElement("jobs"))

		rows, err = database.Connection.Db.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'jobs'")
		Expect(err).NotTo(HaveOccurred())

		defer rows.Close()
		columns := []Column{}

		for rows.Next() {
			var Field, Type string
			err := rows.Scan(&Field, &Type)
			Expect(err).NotTo(HaveOccurred())

			columns = append(columns, Column{
				Field: Field,
				Type:  Type,
			})
		}

		Expect(columns).To(ContainElement(Column{
			Field: "id",
			Type:  "int",
		}))
		Expect(columns).To(ContainElement(Column{
			Field: "worker_id",
			Type:  "varchar",
		}))
		Expect(columns).To(ContainElement(Column{
			Field: "payload",
			Type:  "longtext",
		}))
		Expect(columns).To(ContainElement(Column{
			Field: "version",
			Type:  "bigint",
		}))
		Expect(columns).To(ContainElement(Column{
			Field: "retry_count",
			Type:  "int",
		}))
		Expect(columns).To(ContainElement(Column{
			Field: "active_at",
			Type:  "timestamp",
		}))
	})
})
