package services

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type Registrar struct {
	clientsRepo ClientsRepo
	kindsRepo   KindsRepo
}

func NewRegistrar(clientsRepo ClientsRepo, kindsRepo KindsRepo) Registrar {
	return Registrar{
		clientsRepo: clientsRepo,
		kindsRepo:   kindsRepo,
	}

}

func (registrar Registrar) Register(conn ConnectionInterface, client models.Client, kinds []models.Kind) error {
	_, err := registrar.clientsRepo.Upsert(conn, client)
	if err != nil {
		return err
	}

	for _, kind := range kinds {
		if kind.ID == "" {
			continue
		}

		_, err := registrar.kindsRepo.Upsert(conn, kind)
		if err != nil {
			return err
		}
	}
	return nil
}

func (registrar Registrar) Prune(conn ConnectionInterface, client models.Client, kinds []models.Kind) error {
	kindIDs := []string{}
	for _, kind := range kinds {
		kindIDs = append(kindIDs, kind.ID)
	}

	_, err := registrar.kindsRepo.Trim(conn, client.ID, kindIDs)
	return err
}
