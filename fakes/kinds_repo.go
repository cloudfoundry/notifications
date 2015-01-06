package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

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

func (fake *KindsRepo) Create(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	key := kind.ID + kind.ClientID
	if _, ok := fake.Kinds[key]; ok {
		return kind, models.DuplicateRecordError{}
	}
	fake.Kinds[key] = kind
	return kind, nil
}

func (fake *KindsRepo) Update(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	key := kind.ID + kind.ClientID
	fake.Kinds[key] = kind
	return kind, fake.UpdateError
}

func (fake *KindsRepo) Upsert(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	key := kind.ID + kind.ClientID
	fake.Kinds[key] = kind
	return kind, fake.UpsertError
}

func (fake *KindsRepo) Find(conn models.ConnectionInterface, id, clientID string) (models.Kind, error) {
	key := id + clientID
	if kind, ok := fake.Kinds[key]; ok {
		return kind, fake.FindError
	}
	return models.Kind{}, models.NewRecordNotFoundError("Kind %q %q could not be found", id, clientID)
}

func (fake *KindsRepo) FindByClient(conn models.ConnectionInterface, clientID string) ([]models.Kind, error) {
	kinds := []models.Kind{}

	for _, kind := range fake.Kinds {
		if kind.ClientID == clientID {
			kinds = append(kinds, kind)
		}
	}

	return kinds, nil
}

func (fake *KindsRepo) FindAll(conn models.ConnectionInterface) ([]models.Kind, error) {
	var kinds []models.Kind

	for _, kind := range fake.Kinds {
		kinds = append(kinds, kind)
	}

	return kinds, nil
}

func (fake *KindsRepo) Trim(conn models.ConnectionInterface, clientID string, kindIDs []string) (int, error) {
	fake.TrimArguments = []interface{}{clientID, kindIDs}
	return 0, fake.TrimError
}

func (fake *KindsRepo) FindAllByTemplateID(conn models.ConnectionInterface, templateID string) ([]models.Kind, error) {
	var kinds []models.Kind
	for _, kind := range fake.Kinds {
		if kind.TemplateID == templateID {
			kinds = append(kinds, kind)
		}
	}
	return kinds, fake.FindAllByTemplateIDError
}
