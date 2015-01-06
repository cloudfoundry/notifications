package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateAssociation struct {
	ClientID       string
	NotificationID string
}

type TemplateAssociationListerInterface interface {
	List(string) ([]TemplateAssociation, error)
}

type TemplateAssociationLister struct {
	clientsRepo models.ClientsRepoInterface
	kindsRepo   models.KindsRepoInterface
	database    models.DatabaseInterface
}

func NewTemplateAssociationLister(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface, database models.DatabaseInterface) TemplateAssociationLister {
	return TemplateAssociationLister{
		clientsRepo: clientsRepo,
		kindsRepo:   kindsRepo,
		database:    database,
	}
}

func (lister TemplateAssociationLister) List(templateID string) ([]TemplateAssociation, error) {
	associations := []TemplateAssociation{}

	clients, err := lister.clientsRepo.FindAllByTemplateID(lister.database.Connection(), templateID)
	if err != nil {
		return associations, err
	}

	kinds, err := lister.kindsRepo.FindAllByTemplateID(lister.database.Connection(), templateID)
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
