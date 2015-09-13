package models

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"time"

	sql_migrate "github.com/rubenv/sql-migrate"
)

type DatabaseMigrator struct {
}

func (d DatabaseMigrator) Migrate(db *sql.DB, migrationsPath string) {
	sql_migrate.SetTable("notifications_model_migrations")

	migrations := &sql_migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	_, err := sql_migrate.Exec(db, "mysql", migrations, sql_migrate.Up)
	if err != nil {
		panic(err)
	}
}

func (d DatabaseMigrator) Seed(database DatabaseInterface, defaultTemplatePath string) {
	repo := NewTemplatesRepo()
	bytes, err := ioutil.ReadFile(defaultTemplatePath)
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
		if _, ok := err.(NotFoundError); !ok {
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
