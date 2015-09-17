package collections

import (
	"fmt"

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
	Insert(conn models.ConnectionInterface, template models.Template) (createdTemplate models.Template, err error)
	Update(conn models.ConnectionInterface, template models.Template) (updatedTemplate models.Template, err error)
	Get(conn models.ConnectionInterface, templateID string) (retrievedTemplate models.Template, err error)
	Delete(conn models.ConnectionInterface, templateID string) error
	List(conn models.ConnectionInterface, clientID string) (templateList []models.Template, err error)
}

type TemplatesCollection struct {
	repo templatesRepository
}

func NewTemplatesCollection(repo templatesRepository) TemplatesCollection {
	return TemplatesCollection{
		repo: repo,
	}
}

func (c TemplatesCollection) Set(conn ConnectionInterface, template Template) (Template, error) {
	if template.ID == "" || template.ID == models.DefaultTemplate.ID {
		model, err := c.repo.Insert(conn, models.Template{
			ID:       template.ID,
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
				return c.updateExistingRecord(conn, template)
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

	return c.updateExistingRecord(conn, template)
}

func (c TemplatesCollection) Get(conn ConnectionInterface, templateID, clientID string) (Template, error) {
	template, err := c.repo.Get(conn, templateID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Template{}, NotFoundError{err}
		default:
			return Template{}, PersistenceError{err}
		}
	}

	if template.ClientID != clientID && templateID != "default" {
		return Template{}, NotFoundError{fmt.Errorf("Template with id %q could not be found", templateID)}
	}

	return Template{
		ID:       template.ID,
		Name:     template.Name,
		HTML:     template.HTML,
		Text:     template.Text,
		Subject:  template.Subject,
		Metadata: template.Metadata,
		ClientID: template.ClientID,
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

func (c TemplatesCollection) List(conn ConnectionInterface, clientID string) ([]Template, error) {
	var templateList []Template

	templates, err := c.repo.List(conn, clientID)
	if err != nil {
		return templateList, UnknownError{err}
	}

	for _, template := range templates {
		templateList = append(templateList, Template{
			ID:       template.ID,
			Name:     template.Name,
			HTML:     template.HTML,
			Text:     template.Text,
			Subject:  template.Subject,
			Metadata: template.Metadata,
			ClientID: template.ClientID,
		})
	}

	return templateList, nil
}

func (c TemplatesCollection) updateExistingRecord(conn ConnectionInterface, template Template) (Template, error) {
	model, err := c.repo.Update(conn, models.Template{
		ID:       template.ID,
		Name:     template.Name,
		HTML:     template.HTML,
		Text:     template.Text,
		Subject:  template.Subject,
		Metadata: template.Metadata,
		ClientID: template.ClientID,
	})
	if err != nil {
		return Template{}, PersistenceError{err}
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
