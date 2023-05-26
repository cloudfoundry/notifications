package middleware_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database Allocator", func() {
	var (
		ware    middleware.DatabaseAllocator
		sqlDB   *sql.DB
		writer  *httptest.ResponseRecorder
		request *http.Request
		buffer  *bytes.Buffer
		context stack.Context
	)

	BeforeEach(func() {
		var err error
		sqlDB, _, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		ware = middleware.NewDatabaseAllocator(sqlDB, true)

		writer = httptest.NewRecorder()
		request = &http.Request{}

		buffer = bytes.NewBuffer([]byte{})
		logger := lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))

		context = stack.NewContext()
		context.Set("logger", logger)
		context.Set(middleware.VCAPRequestIDKey, "some-vcap-request-id")
	})

	It("allocates a database that it adds to the context", func() {
		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		database, ok := context.Get("database").(*db.DB)
		Expect(ok).To(BeTrue())

		connection, ok := database.Connection().(*db.Connection)
		Expect(ok).To(BeTrue())
		Expect(connection.DbMap.Db).To(Equal(sqlDB))

		_, err := connection.DbMap.TableFor(reflect.TypeOf(models.Client{}), false)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when db logging is enabled", func() {
		It("allocates a database that logs with the VCAP request ID", func() {
			result := ware.ServeHTTP(writer, request, context)
			Expect(result).To(BeTrue())

			database := context.Get("database").(*db.DB)
			conn := database.Connection()
			conn.Exec("SELECT * FROM `posts` WHERE `id` = ? AND `tag` = ?", 1234, "banana")

			logs := buffer.Bytes()
			Expect(len(logs)).To(BeNumerically(">", 0))

			var line logLine
			err := json.Unmarshal(logs, &line)
			Expect(err).NotTo(HaveOccurred())
			Expect(line.Source).To(Equal("notifications"))
			Expect(line.Message).To(Equal("notifications.db"))
			Expect(line.LogLevel).To(Equal(int(lager.DEBUG)))
			Expect(line.Data).To(HaveKeyWithValue("vcap_request_id", "some-vcap-request-id"))
			Expect(line.Data).To(HaveKeyWithValue("statement", MatchRegexp(`^SELECT \* FROM `+"`"+`posts`+"`"+` WHERE `+"`"+`id`+"`"+` = \? AND `+"`"+`tag`+"`"+` = \? \[1:1234 2:"banana"\]$`)))
		})
	})

	Context("when db logging is disabled", func() {
		It("allocates a database that does not log", func() {
			ware = middleware.NewDatabaseAllocator(sqlDB, false)

			result := ware.ServeHTTP(writer, request, context)
			Expect(result).To(BeTrue())

			database := context.Get("database").(*db.DB)
			conn := database.Connection()
			conn.Exec("SELECT * FROM `posts` WHERE `id` = ? AND `tag` = ?", 1234, "banana")

			Expect(buffer.Bytes()).To(HaveLen(0))
		})
	})
})
