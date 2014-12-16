package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateAssignerInterface interface {
	AssignToClient(string, string) error
}

type TemplateAssigner struct {
	clientsRepo   models.ClientsRepoInterface
	templatesRepo models.TemplatesRepoInterface
	database      models.DatabaseInterface
}

func NewTemplateAssigner(clientsRepo models.ClientsRepoInterface,
	templatesRepo models.TemplatesRepoInterface,
	database models.DatabaseInterface) TemplateAssigner {
	return TemplateAssigner{
		clientsRepo:   clientsRepo,
		templatesRepo: templatesRepo,
		database:      database,
	}
}

func (assigner TemplateAssigner) AssignToClient(clientID, templateID string) error {
	conn := assigner.database.Connection()

	client, err := assigner.clientsRepo.Find(conn, clientID)
	if err != nil {
		if (err == models.ErrRecordNotFound{}) {
			return ClientMissingError("No client with id '" + clientID + "'")
		}
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
	return err
}
