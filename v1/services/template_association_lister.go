package services

type TemplateAssociation struct {
	ClientID       string
	NotificationID string
}

type TemplateAssociationLister struct {
	clientsRepo   ClientsRepo
	kindsRepo     KindsRepo
	templatesRepo TemplatesRepo
}

func NewTemplateAssociationLister(clientsRepo ClientsRepo, kindsRepo KindsRepo, templatesRepo TemplatesRepo) TemplateAssociationLister {
	return TemplateAssociationLister{
		clientsRepo:   clientsRepo,
		kindsRepo:     kindsRepo,
		templatesRepo: templatesRepo,
	}
}

func (lister TemplateAssociationLister) List(database DatabaseInterface, templateID string) ([]TemplateAssociation, error) {
	associations := []TemplateAssociation{}
	conn := database.Connection()

	_, err := lister.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		return associations, err
	}

	clients, err := lister.clientsRepo.FindAllByTemplateID(conn, templateID)
	if err != nil {
		return associations, err
	}

	kinds, err := lister.kindsRepo.FindAllByTemplateID(conn, templateID)
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
