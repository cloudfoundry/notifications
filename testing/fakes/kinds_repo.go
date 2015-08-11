package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type KindsRepo struct {
	Kinds                    map[string]models.Kind
	UpsertError              error
	TrimError                error
	FindError                error
	UpdateError              error
	FindAllByTemplateIDError error
	TrimArguments            []interface{}
}

func NewKindsRepo() *KindsRepo {
	return &KindsRepo{
		Kinds:         make(map[string]models.Kind),
		TrimArguments: make([]interface{}, 0),
	}
}

func (fake *KindsRepo) Create(conn db.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	if kind.TemplateID == "" {
		kind.TemplateID = models.DefaultTemplateID
	}

	key := kind.ID + kind.ClientID
	if _, ok := fake.Kinds[key]; ok {
		return kind, models.DuplicateRecordError{}
	}

	fake.Kinds[key] = kind
	return kind, nil
}

func (fake *KindsRepo) Update(conn db.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	if fake.UpdateError != nil {
		return kind, fake.UpdateError
	}

	if kind.TemplateID == "" {
		existingKind, err := fake.Find(conn, kind.ID, kind.ClientID)
		if err != nil {
			return kind, err
		}
		kind.TemplateID = existingKind.TemplateID
	}

	key := kind.ID + kind.ClientID
	fake.Kinds[key] = kind
	return kind, nil
}

func (fake *KindsRepo) Upsert(conn db.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	key := kind.ID + kind.ClientID
	fake.Kinds[key] = kind
	return kind, fake.UpsertError
}

func (fake *KindsRepo) Find(conn db.ConnectionInterface, id, clientID string) (models.Kind, error) {
	if fake.FindError != nil {
		return models.Kind{}, fake.FindError
	}
	key := id + clientID
	if kind, ok := fake.Kinds[key]; ok {
		return kind, nil
	}
	return models.Kind{}, models.NewRecordNotFoundError("Kind %q %q could not be found", id, clientID)
}

func (fake *KindsRepo) FindByClient(conn db.ConnectionInterface, clientID string) ([]models.Kind, error) {
	kinds := []models.Kind{}

	for _, kind := range fake.Kinds {
		if kind.ClientID == clientID {
			kinds = append(kinds, kind)
		}
	}

	return kinds, nil
}

func (fake *KindsRepo) FindAll(conn db.ConnectionInterface) ([]models.Kind, error) {
	var kinds []models.Kind

	for _, kind := range fake.Kinds {
		kinds = append(kinds, kind)
	}

	return kinds, nil
}

func (fake *KindsRepo) Trim(conn db.ConnectionInterface, clientID string, kindIDs []string) (int, error) {
	fake.TrimArguments = []interface{}{clientID, kindIDs}
	return 0, fake.TrimError
}

func (fake *KindsRepo) FindAllByTemplateID(conn db.ConnectionInterface, templateID string) ([]models.Kind, error) {
	var kinds []models.Kind
	for _, kind := range fake.Kinds {
		if kind.TemplateID == templateID {
			kinds = append(kinds, kind)
		}
	}
	return kinds, fake.FindAllByTemplateIDError
}
