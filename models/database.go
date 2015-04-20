package models

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"time"

	"sync"

	"github.com/coopernurse/gorp"
	sql_migrate "github.com/rubenv/sql-migrate"

	_ "github.com/go-sql-driver/mysql"
)

var _database *DB
var mutex sync.Mutex

type DB struct {
	connection *Connection
	config     Config
}

type DatabaseInterface interface {
	Connection() ConnectionInterface
	TraceOn(string, gorp.GorpLogger)
	Seed()
}

type configurableDB interface {
	SetMaxOpenConns(int)
}

func NewDatabase(config Config) *DB {
	if _database != nil {
		return _database
	}

	mutex.Lock()
	defer mutex.Unlock()
	db, err := sql.Open("mysql", config.DatabaseURL)
	if err != nil {
		panic(err)
	}

	ConfigureDB(db, config)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	connection := &Connection{
		DbMap: &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Engine:   "InnoDB",
				Encoding: "UTF8",
			},
		},
	}

	_database = &DB{
		config:     config,
		connection: connection,
	}

	_database.migrate(config.MigrationsPath)

	return _database
}

func ConfigureDB(db configurableDB, config Config) {
	db.SetMaxOpenConns(config.MaxOpenConnections)
}

func (database DB) migrate(migrationsPath string) {
	sql_migrate.SetTable("notifications_model_migrations")

	migrations := &sql_migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	_, err := sql_migrate.Exec(database.connection.Db, "mysql", migrations, sql_migrate.Up)
	if err != nil {
		panic(err)
	}

	database.connection.AddTableWithName(Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
	database.connection.AddTableWithName(Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
	database.connection.AddTableWithName(Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")
	database.connection.AddTableWithName(Unsubscribe{}, "unsubscribes").SetKeys(true, "Primary").SetUniqueTogether("user_id", "client_id", "kind_id")
	database.connection.AddTableWithName(GlobalUnsubscribe{}, "global_unsubscribes").SetKeys(true, "Primary").ColMap("UserID").SetUnique(true)
	database.connection.AddTableWithName(Template{}, "templates").SetKeys(true, "Primary").ColMap("Name").SetUnique(true)
	database.connection.AddTableWithName(Message{}, "messages").SetKeys(false, "ID")
}

func (database DB) Seed() {
	repo := NewTemplatesRepo()
	bytes, err := ioutil.ReadFile(database.config.DefaultTemplatePath)
	if err != nil {
		panic(err)
	}

	var template struct {
		Name     string          `json:"name"`
		Subject  string          `json:"subject"`
		Text     string          `json:"text"`
		HTML     string          `json:"html"`
		Metadata json.RawMessage `json:"metadata"`
	}

	err = json.Unmarshal(bytes, &template)
	if err != nil {
		panic(err)
	}

	conn := database.Connection()
	existingTemplate, err := repo.FindByID(conn, DefaultTemplateID)
	if err != nil {
		if _, ok := err.(RecordNotFoundError); !ok {
			panic(err)
		}

		_, err = repo.Create(conn, Template{
			ID:       DefaultTemplateID,
			Name:     template.Name,
			Subject:  template.Subject,
			HTML:     template.HTML,
			Text:     template.Text,
			Metadata: string(template.Metadata),
		})
		if err != nil {
			panic(err)
		}

		return
	}

	if !existingTemplate.Overridden {
		existingTemplate.Name = template.Name
		existingTemplate.Subject = template.Subject
		existingTemplate.HTML = template.HTML
		existingTemplate.Text = template.Text
		existingTemplate.Metadata = string(template.Metadata)
		existingTemplate.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()
		_, err = conn.Update(&existingTemplate)
		if err != nil {
			panic(err)
		}
	}
}

func (database *DB) Connection() ConnectionInterface {
	return database.connection
}

func (database *DB) TraceOn(prefix string, logger gorp.GorpLogger) {
	database.connection.TraceOn(prefix, logger)
}
