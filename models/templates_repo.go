package models

import "database/sql"

type TemplatesRepoInterface interface {
    Find(ConnectionInterface, string) (Template, error)
    Upsert(ConnectionInterface, Template) (Template, error)
}

type TemplatesRepo struct{}

func NewTemplatesRepo() TemplatesRepo {
    return TemplatesRepo{}
}

func (repo TemplatesRepo) Find(conn ConnectionInterface, templateName string) (Template, error) {
    template := Template{}
    err := conn.SelectOne(&template, "SELECT * FROM `templates` WHERE `name`=?", templateName)
    if err != nil {
        if err == sql.ErrNoRows {
            return template, ErrRecordNotFound{}
        }
        return template, err
    }
    return template, nil
}

func (repo TemplatesRepo) Upsert(conn ConnectionInterface, template Template) (Template, error) {
    var err error
    _, err = repo.Find(conn, template.Name)
    if err != nil {
        if (err == ErrRecordNotFound{}) {
            err = conn.Insert(&template)
        }
    }
    _, err = conn.Update(&template)

    return template, err
}
