package middleware

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"
)

type DatabaseAllocator struct {
	DB    *sql.DB
	Trace bool
}

func NewDatabaseAllocator(sqlDB *sql.DB, trace bool) DatabaseAllocator {
	return DatabaseAllocator{
		DB:    sqlDB,
		Trace: trace,
	}
}

func (ware DatabaseAllocator) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) bool {
	db := models.NewDatabase(ware.DB, models.Config{})

	logger := gorpCompatibleLogger{
		logger: context.Get("logger").(lager.Logger).WithData(lager.Data{
			VCAPRequestIDKey: context.Get(VCAPRequestIDKey),
		}),
	}

	if ware.Trace {
		db.TraceOn("", logger)
	}

	context.Set("database", db)
	return true
}

type gorpCompatibleLogger struct {
	logger lager.Logger
}

func (g gorpCompatibleLogger) Printf(format string, v ...interface{}) {
	g.logger.Debug("db", lager.Data{
		"statement": fmt.Sprintf(format, v...),
	})
}
