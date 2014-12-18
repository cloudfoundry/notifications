package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateAssignerInterface interface {
	AssignToClient(string, string) error
	AssignToNotification(string, string, string) error
}

type TemplateAssigner struct {
	clientsRepo   models.ClientsRepoInterface
	kindsRepo     models.KindsRepoInterface
	templatesRepo models.TemplatesRepoInterface
	database      models.DatabaseInterface
}

func NewTemplateAssigner(clientsRepo models.ClientsRepoInterface,
	kindsRepo models.KindsRepoInterface,
	templatesRepo models.TemplatesRepoInterface,
	database models.DatabaseInterface) TemplateAssigner {

	return TemplateAssigner{
		clientsRepo:   clientsRepo,
		kindsRepo:     kindsRepo,
		templatesRepo: templatesRepo,
		database:      database,
	}
}

func (assigner TemplateAssigner) AssignToClient(clientID, templateID string) error {
	conn := assigner.database.Connection()

	client, err := assigner.clientsRepo.Find(conn, clientID)
	if err != nil {
		return err
	}

	_, err = assigner.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		if (err == models.ErrRecordNotFound{}) {
			return TemplateAssignmentError("No template with id '" + templateID + "'")
		}
		return err
	}

	client.Template = templateID

	_, err = assigner.clientsRepo.Update(conn, client)
	if err != nil {
		return err
	}

	return nil
}

func (assigner TemplateAssigner) AssignToNotification(clientID, notificationID, templateID string) error {
	conn := assigner.database.Connection()

	_, err := assigner.clientsRepo.Find(conn, clientID)
	if err != nil {
		return err
	}

	kind, err := assigner.kindsRepo.Find(conn, notificationID, clientID)
	if err != nil {
		return err
	}

	_, err = assigner.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		if (err == models.ErrRecordNotFound{}) {
			return TemplateAssignmentError("No template with id '" + templateID + "'")
		}
		return err
	}

	kind.Template = templateID
	_, err = assigner.kindsRepo.Update(conn, kind)
	if err != nil {
		return err
	}

	return nil
}
