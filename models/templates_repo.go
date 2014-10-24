package models

import "database/sql"

type TemplatesRepoInterface interface {
    Find(ConnectionInterface, string) (Template, error)
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
