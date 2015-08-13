package models

import (
	"database/sql"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/db"
)

type TemplatesRepository struct {
	generateGUID guidGeneratorFunc
}

func NewTemplatesRepository(guidGenerator guidGeneratorFunc) TemplatesRepository {
	return TemplatesRepository{
		generateGUID: guidGenerator,
	}
}

func (r TemplatesRepository) Insert(conn db.ConnectionInterface, template Template) (Template, error) {
	guid, err := r.generateGUID()
	if err != nil {
		panic(err)
	}

	present, err := r.templateWithNameAndClientIDIsPresent(conn, template.Name, template.ClientID)
	if err != nil {
		return template, err
	}
	if present {
		return template, DuplicateRecordError{fmt.Errorf("Template with name %q already exists", template.Name)}
	}

	template.ID = guid.String()
	err = conn.Insert(&template)
	if err != nil {
		return template, err
	}

	return template, nil
}

func (r TemplatesRepository) Get(conn db.ConnectionInterface, templateID string) (Template, error) {
	template := Template{}
	err := conn.SelectOne(&template, "SELECT * FROM `v2_templates` WHERE `id` = ?", templateID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = RecordNotFoundError{fmt.Errorf("Template with id %q could not be found", templateID)}
		}

		return template, err
	}
	return template, nil
}

func (r TemplatesRepository) templateWithNameAndClientIDIsPresent(conn db.ConnectionInterface, name, clientID string) (bool, error) {
	err := conn.SelectOne(&Template{}, "SELECT * FROM `v2_templates` WHERE `name` = ? AND `client_id` = ?", name, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
