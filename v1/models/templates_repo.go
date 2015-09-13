package models

import (
	"database/sql"
	"fmt"
	"time"
)

type TemplatesRepo struct{}

func NewTemplatesRepo() TemplatesRepo {
	return TemplatesRepo{}
}

func (repo TemplatesRepo) FindByID(conn ConnectionInterface, templateID string) (Template, error) {
	template := Template{}
	err := conn.SelectOne(&template, "SELECT * FROM `templates` WHERE `id`=?", templateID)
	if err != nil {
		if err == sql.ErrNoRows {
			return template, NotFoundError{fmt.Errorf("Template with ID %q could not be found", templateID)}
		}
		return template, err
	}
	return template, nil
}

func (repo TemplatesRepo) Update(conn ConnectionInterface, templateID string, template Template) (Template, error) {
	existingTemplate, err := repo.FindByID(conn, templateID)
	if err != nil {
		return existingTemplate, err
	}

	template.Primary = existingTemplate.Primary
	template.ID = existingTemplate.ID
	template.CreatedAt = existingTemplate.CreatedAt
	template.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()
	template.Overridden = true

	_, err = conn.Update(&template)
	if err != nil {
		return Template{}, TemplateUpdateError{err}
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
	err := conn.Insert(&template)
	if err != nil {
		return Template{}, err
	}

	return template, nil
}

func (repo TemplatesRepo) Destroy(conn ConnectionInterface, templateID string) error {
	template, err := repo.FindByID(conn, templateID)
	if err != nil {
		return err
	}

	_, err = conn.Delete(&template)

	return err
}
