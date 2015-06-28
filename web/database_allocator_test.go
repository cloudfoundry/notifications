package web_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database Allocator", func() {
	var (
		ware    web.DatabaseAllocator
		sqlDB   *sql.DB
		writer  *httptest.ResponseRecorder
		request *http.Request
		buffer  *bytes.Buffer
		context stack.Context
	)

	BeforeEach(func() {
		var err error
		sqlDB, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		ware = web.NewDatabaseAllocator(sqlDB, true)

		writer = httptest.NewRecorder()
		request = &http.Request{}

		buffer = bytes.NewBuffer([]byte{})
		logger := lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))

		context = stack.NewContext()
		context.Set("logger", logger)
		context.Set(handlers.VCAPRequestIDKey, "some-vcap-request-id")
	})

	AfterEach(func() {
		err := sqlDB.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("allocates a database that it adds to the context", func() {
		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		database, ok := context.Get("database").(*models.DB)
		Expect(ok).To(BeTrue())

		connection, ok := database.Connection().(*models.Connection)
		Expect(ok).To(BeTrue())
		Expect(connection.DbMap.Db).To(Equal(sqlDB))

		_, err := connection.DbMap.TableFor(reflect.TypeOf(models.Client{}), false)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when db logging is enabled", func() {
		It("allocates a database that logs with the VCAP request ID", func() {
			result := ware.ServeHTTP(writer, request, context)
			Expect(result).To(BeTrue())

			database := context.Get("database").(*models.DB)
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
			Expect(line.Data).To(HaveKeyWithValue("statement", "SELECT * FROM `posts` WHERE `id` = ? AND `tag` = ? [1:1234 2:\"banana\"]"))
		})
	})

	Context("when db logging is disabled", func() {
		It("allocates a database that does not log", func() {
			ware = web.NewDatabaseAllocator(sqlDB, false)

			result := ware.ServeHTTP(writer, request, context)
			Expect(result).To(BeTrue())

			database := context.Get("database").(*models.DB)
			conn := database.Connection()
			conn.Exec("SELECT * FROM `posts` WHERE `id` = ? AND `tag` = ?", 1234, "banana")

			Expect(buffer.Bytes()).To(HaveLen(0))
		})
	})
})
