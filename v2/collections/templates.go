package collections

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type Template struct {
	ID       string
	Name     string
	HTML     string
	Text     string
	Subject  string
	Metadata string
	ClientID string
}

type templatesRepository interface {
	Insert(conn db.ConnectionInterface, template models.Template) (createdTemplate models.Template, err error)
	Get(conn db.ConnectionInterface, templateID string) (retrievedTemplate models.Template, err error)
	Delete(conn db.ConnectionInterface, templateID string) error
}

type TemplatesCollection struct {
	repo templatesRepository
}

func NewTemplatesCollection(repo templatesRepository) TemplatesCollection {
	return TemplatesCollection{
		repo: repo,
	}
}

func (c TemplatesCollection) Set(conn ConnectionInterface, template Template) (createdTemplate Template, err error) {
	model, err := c.repo.Insert(conn, models.Template{
		Name:     template.Name,
		HTML:     template.HTML,
		Text:     template.Text,
		Subject:  template.Subject,
		Metadata: template.Metadata,
		ClientID: template.ClientID,
	})
	if err != nil {
		switch err.(type) {
		case models.DuplicateRecordError:
			return Template{}, DuplicateRecordError{err}
		default:
			return Template{}, PersistenceError{err}
		}
	}

	return Template{
		ID:       model.ID,
		Name:     model.Name,
		HTML:     model.HTML,
		Text:     model.Text,
		Subject:  model.Subject,
		Metadata: model.Metadata,
		ClientID: model.ClientID,
	}, nil
}

func (c TemplatesCollection) Get(conn ConnectionInterface, templateID, clientID string) (Template, error) {
	model, err := c.repo.Get(conn, templateID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Template{}, NotFoundError{err}
		default:
			return Template{}, PersistenceError{err}
		}
	}
	if model.ClientID != clientID {
		return Template{}, NotFoundError{fmt.Errorf("Template with id %q could not be found", templateID)}
	}

	return Template{
		ID:       model.ID,
		Name:     model.Name,
		HTML:     model.HTML,
		Text:     model.Text,
		Subject:  model.Subject,
		Metadata: model.Metadata,
		ClientID: model.ClientID,
	}, nil
}

func (c TemplatesCollection) Delete(conn ConnectionInterface, templateID string) error {
	err := c.repo.Delete(conn, templateID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return NotFoundError{err}
		default:
			return PersistenceError{err}
		}
	}

	return nil
}
