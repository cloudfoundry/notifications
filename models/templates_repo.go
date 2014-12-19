package models

import (
	"database/sql"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type TemplatesRepoInterface interface {
	FindByID(ConnectionInterface, string) (Template, error)
	Find(ConnectionInterface, string) (Template, error)
	Create(ConnectionInterface, Template) (Template, error)
	Update(ConnectionInterface, string, Template) (Template, error)
	Upsert(ConnectionInterface, Template) (Template, error)
	ListIDsAndNames(ConnectionInterface) ([]Template, error)
	Destroy(ConnectionInterface, string) error
	DeprecatedDestroy(ConnectionInterface, string) error
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

func (repo TemplatesRepo) Find(conn ConnectionInterface, templateName string) (Template, error) {
	template := Template{}
	err := conn.SelectOne(&template, "SELECT * FROM `templates` WHERE `name`=?", templateName)
	if err != nil {
		if err == sql.ErrNoRows {
			return template, NewRecordNotFoundError("Template with name %q could not be found", templateName)
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

	existingTemplate.Name = template.Name
	existingTemplate.Subject = template.Subject
	existingTemplate.HTML = template.HTML
	existingTemplate.Text = template.Text
	existingTemplate.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()

	_, err = conn.Update(&existingTemplate)
	if err != nil {
		return Template{}, TemplateUpdateError{Message: err.Error()}
	}

	return template, nil
}

func (repo TemplatesRepo) Upsert(conn ConnectionInterface, template Template) (Template, error) {
	existingTemplate, err := repo.Find(conn, template.Name)
	if err != nil {
		if _, ok := err.(RecordNotFoundError); ok {
			return repo.Create(conn, template)
		}
		return Template{}, err
	}

	template.Primary = existingTemplate.Primary
	template.CreatedAt = existingTemplate.CreatedAt
	template.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()
	_, err = conn.Update(&template)
	if err != nil {
		return Template{}, err
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
	setTemplateTimestamps(&template)
	template.ID = uuid.New()

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

func (repo TemplatesRepo) DeprecatedDestroy(conn ConnectionInterface, templateName string) error {
	template, err := repo.Find(conn, templateName)
	if err != nil {
		if _, ok := err.(RecordNotFoundError); ok {
			return nil
		}
		return err
	}

	_, err = conn.Delete(&template)

	return err
}
