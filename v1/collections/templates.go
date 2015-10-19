package collections

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type DatabaseInterface interface {
	models.DatabaseInterface
}

type ConnectionInterface interface {
	models.ConnectionInterface
}

type TemplateAssignmentError struct {
	Err error
}

func (e TemplateAssignmentError) Error() string {
	return e.Err.Error()
}

type clientsRepository interface {
	Find(connection models.ConnectionInterface, clientID string) (models.Client, error)
	FindAllByTemplateID(connection models.ConnectionInterface, templateID string) ([]models.Client, error)
	Update(connection models.ConnectionInterface, client models.Client) (models.Client, error)
}

type kindsRepository interface {
	Find(connection models.ConnectionInterface, kindID string, clientID string) (models.Kind, error)
	FindAllByTemplateID(connection models.ConnectionInterface, templateID string) ([]models.Kind, error)
	Update(connection models.ConnectionInterface, kind models.Kind) (models.Kind, error)
}

type templatesRepository interface {
	FindByID(connection models.ConnectionInterface, templateID string) (models.Template, error)
	Create(connection models.ConnectionInterface, template models.Template) (models.Template, error)
	Destroy(connection models.ConnectionInterface, templateID string) error
}

type TemplateAssociation struct {
	ClientID       string
	NotificationID string
}

type Template struct {
	ID       string
	Name     string
	Text     string
	HTML     string
	Subject  string
	Metadata string
}

type TemplatesCollection struct {
	clientsRepo   clientsRepository
	kindsRepo     kindsRepository
	templatesRepo templatesRepository
}

func NewTemplatesCollection(clientsRepo clientsRepository, kindsRepo kindsRepository, templatesRepo templatesRepository) TemplatesCollection {
	return TemplatesCollection{
		clientsRepo:   clientsRepo,
		kindsRepo:     kindsRepo,
		templatesRepo: templatesRepo,
	}
}

func (c TemplatesCollection) AssignToClient(conn ConnectionInterface, clientID, templateID string) error {
	if templateID == "" {
		templateID = models.DefaultTemplateID
	}

	client, err := c.clientsRepo.Find(conn, clientID)
	if err != nil {
		return err
	}

	err = c.findTemplate(conn, templateID)
	if err != nil {
		return err
	}

	client.TemplateID = templateID

	_, err = c.clientsRepo.Update(conn, client)
	if err != nil {
		return err
	}

	return nil
}

func (c TemplatesCollection) AssignToNotification(conn ConnectionInterface, clientID, notificationID, templateID string) error {
	if templateID == "" {
		templateID = models.DefaultTemplateID
	}

	_, err := c.clientsRepo.Find(conn, clientID)
	if err != nil {
		return err
	}

	kind, err := c.kindsRepo.Find(conn, notificationID, clientID)
	if err != nil {
		return err
	}

	err = c.findTemplate(conn, templateID)
	if err != nil {
		return err
	}

	kind.TemplateID = templateID

	_, err = c.kindsRepo.Update(conn, kind)
	if err != nil {
		return err
	}

	return nil
}

func (c TemplatesCollection) findTemplate(conn ConnectionInterface, templateID string) error {
	if templateID == "" {
		return nil
	}

	_, err := c.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		if _, ok := err.(models.NotFoundError); ok {
			return TemplateAssignmentError{fmt.Errorf("No template with id %q", templateID)}
		}
		return err
	}

	return nil
}

func (c TemplatesCollection) ListAssociations(conn ConnectionInterface, templateID string) ([]TemplateAssociation, error) {
	associations := []TemplateAssociation{}

	_, err := c.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		return associations, err
	}

	clients, err := c.clientsRepo.FindAllByTemplateID(conn, templateID)
	if err != nil {
		return associations, err
	}

	kinds, err := c.kindsRepo.FindAllByTemplateID(conn, templateID)
	if err != nil {
		return associations, err
	}

	for _, client := range clients {
		associations = append(associations, TemplateAssociation{
			ClientID: client.ID,
		})
	}

	for _, kind := range kinds {
		associations = append(associations, TemplateAssociation{
			ClientID:       kind.ClientID,
			NotificationID: kind.ID,
		})
	}

	return associations, nil
}

func (c TemplatesCollection) Create(connection ConnectionInterface, template Template) (Template, error) {
	tmpl, err := c.templatesRepo.Create(connection, models.Template{
		Name:     template.Name,
		Text:     template.Text,
		HTML:     template.HTML,
		Subject:  template.Subject,
		Metadata: template.Metadata,
	})
	if err != nil {
		return Template{}, err
	}

	return Template{
		ID:       tmpl.ID,
		Name:     tmpl.Name,
		Text:     tmpl.Text,
		HTML:     tmpl.HTML,
		Subject:  tmpl.Subject,
		Metadata: tmpl.Metadata,
	}, nil
}

func (c TemplatesCollection) Delete(connection ConnectionInterface, templateID string) error {
	return c.templatesRepo.Destroy(connection, templateID)
}
