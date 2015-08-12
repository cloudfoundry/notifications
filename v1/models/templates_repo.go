package models

import (
	"database/sql"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
)

type TemplatesRepo struct{}

func NewTemplatesRepo() TemplatesRepo {
	return TemplatesRepo{}
}

func (repo TemplatesRepo) FindByID(conn db.ConnectionInterface, templateID string) (Template, error) {
	template := Template{}
	err := conn.SelectOne(&template, "SELECT * FROM `templates` WHERE `id`=?", templateID)
	if err != nil {
		if err == sql.ErrNoRows {
			return template, NewRecordNotFoundError("Template with ID %q could not be found", templateID)
		}
		return template, err
	}
	return template, nil
}

func (repo TemplatesRepo) Update(conn db.ConnectionInterface, templateID string, template Template) (Template, error) {
	existingTemplate, err := repo.FindByID(conn, templateID)
	if err != nil {
		if _, ok := err.(RecordNotFoundError); ok {
			return existingTemplate, TemplateFindError{Message: "Template " + templateID + " not found"}
		}
		return existingTemplate, err
	}

	template.Primary = existingTemplate.Primary
	template.ID = existingTemplate.ID
	template.CreatedAt = existingTemplate.CreatedAt
	template.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()
	template.Overridden = true

	_, err = conn.Update(&template)
	if err != nil {
		return Template{}, TemplateUpdateError{Message: err.Error()}
	}

	return template, nil
}

func (repo TemplatesRepo) ListIDsAndNames(conn db.ConnectionInterface) ([]Template, error) {
	templates := []Template{}
	_, err := conn.Select(&templates, "SELECT ID, Name FROM `templates`")
	if err != nil {
		return []Template{}, err
	}
	return templates, nil
}

func (repo TemplatesRepo) Create(conn db.ConnectionInterface, template Template) (Template, error) {
	err := conn.Insert(&template)
	if err != nil {
		return Template{}, err
	}

	return template, nil
}

func (repo TemplatesRepo) Destroy(conn db.ConnectionInterface, templateID string) error {
	template, err := repo.FindByID(conn, templateID)
	if err != nil {
		return err
	}

	_, err = conn.Delete(&template)

	return err
}
