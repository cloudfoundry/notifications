package gobble_test

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo"
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
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		tables := []string{}

		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				panic(err)
			}
			tables = append(tables, table)
		}

		Expect(tables).To(ContainElement("jobs"))

		rows, err = database.Connection.Db.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'jobs'")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		columns := []Column{}

		for rows.Next() {
			var Field, Type string
			if err := rows.Scan(&Field, &Type); err != nil {
				panic(err)
			}
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
