package services

import "github.com/cloudfoundry-incubator/notifications/models"

type RegistrarInterface interface {
    Register(models.ConnectionInterface, models.Client, []models.Kind) error
    Prune(models.ConnectionInterface, models.Client, []models.Kind) error
}

type Registrar struct {
    clientsRepo models.ClientsRepoInterface
    kindsRepo   models.KindsRepoInterface
}

func NewRegistrar(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface) Registrar {
    return Registrar{
        clientsRepo: clientsRepo,
        kindsRepo:   kindsRepo,
    }

}

func (registrar Registrar) Register(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    client, err := registrar.clientsRepo.Upsert(conn, client)
    if err != nil {
        return err
    }

    for _, kind := range kinds {
        _, err := registrar.kindsRepo.Upsert(conn, kind)
        if err != nil {
            return err
        }
    }
    return nil
}

func (registrar Registrar) Prune(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    kindIDs := []string{}
    for _, kind := range kinds {
        kindIDs = append(kindIDs, kind.ID)
    }

    _, err := registrar.kindsRepo.Trim(conn, client.ID, kindIDs)
    return err
}
