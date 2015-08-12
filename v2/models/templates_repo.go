package models

import "github.com/cloudfoundry-incubator/notifications/db"

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

	template.ID = guid.String()
	err = conn.Insert(&template)
	if err != nil {
		panic(err)
	}

	return template, nil
}

func (r TemplatesRepository) Get(conn db.ConnectionInterface, templateID string) (Template, error) {
	template := Template{}
	err := conn.SelectOne(&template, "SELECT * FROM `v2_templates` WHERE `id` = ?", templateID)
	if err != nil {
		panic(err)
	}
	return template, nil
}
