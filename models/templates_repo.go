package models

import (
	"database/sql"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type TemplatesRepoInterface interface {
	FindByID(ConnectionInterface, string) (Template, error)
	Create(ConnectionInterface, Template) (Template, error)
	Update(ConnectionInterface, string, Template) (Template, error)
	ListIDsAndNames(ConnectionInterface) ([]Template, error)
	Destroy(ConnectionInterface, string) error
}

type TemplatesRepo struct{}

func NewTemplatesRepo() TemplatesRepo {
	return TemplatesRepo{}
}

func (repo TemplatesRepo) FindByID(conn ConnectionInterface, templateID string) (Template, error) {
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

func (repo TemplatesRepo) Update(conn ConnectionInterface, templateID string, template Template) (Template, error) {
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

func (repo TemplatesRepo) ListIDsAndNames(conn ConnectionInterface) ([]Template, error) {
	templates := []Template{}
	_, err := conn.Select(&templates, "SELECT ID, Name FROM `templates`")
	if err != nil {
		return []Template{}, err
	}
	return templates, nil
}

func (repo TemplatesRepo) Create(conn ConnectionInterface, template Template) (Template, error) {
	template.ID = uuid.New()

	return repo.create(conn, template)
}

func (repo TemplatesRepo) create(conn ConnectionInterface, template Template) (Template, error) {
	setTemplateTimestamps(&template)
	err := conn.Insert(&template)
	if err != nil {
		return Template{}, err
	}

	return template, nil
}

func setTemplateTimestamps(template *Template) {
	if (template.CreatedAt == time.Time{}) {
		template.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
	}
	template.UpdatedAt = template.CreatedAt
}

func (repo TemplatesRepo) Destroy(conn ConnectionInterface, templateID string) error {
	template, err := repo.FindByID(conn, templateID)
	if err != nil {
		return err
	}

	_, err = conn.Delete(&template)

	return err
}
