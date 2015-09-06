package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type KindsRepo struct {
	FindCall struct {
		CallCount int
		Receives  struct {
			Connection models.ConnectionInterface
			KindID     string
			ClientID   string
		}
		Returns struct {
			Kinds []models.Kind
			Error error
		}
	}

	FindAllCall struct {
		Receives struct {
			Connection models.ConnectionInterface
		}
		Returns struct {
			Kinds []models.Kind
			Error error
		}
	}

	FindAllByTemplateIDCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Kinds []models.Kind
			Error error
		}
	}

	TrimCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			ClientID   string
			KindIDs    []string
		}
		Returns struct {
			AffectedRowCount int
			Error            error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Kind       models.Kind
		}
		Returns struct {
			Kind  models.Kind
			Error error
		}
	}

	UpsertCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Kinds      []models.Kind
		}
		Returns struct {
			Kind  models.Kind
			Error error
		}
	}
}

func NewKindsRepo() *KindsRepo {
	return &KindsRepo{}
}

func (kr *KindsRepo) Find(conn models.ConnectionInterface, kindID, clientID string) (models.Kind, error) {
	kr.FindCall.Receives.Connection = conn
	kr.FindCall.Receives.KindID = kindID
	kr.FindCall.Receives.ClientID = clientID

	kind := kr.FindCall.Returns.Kinds[kr.FindCall.CallCount]
	kr.FindCall.CallCount++

	return kind, kr.FindCall.Returns.Error
}

func (kr *KindsRepo) FindAll(conn models.ConnectionInterface) ([]models.Kind, error) {
	kr.FindAllCall.Receives.Connection = conn

	return kr.FindAllCall.Returns.Kinds, kr.FindAllCall.Returns.Error
}

func (kr *KindsRepo) FindAllByTemplateID(conn models.ConnectionInterface, templateID string) ([]models.Kind, error) {
	kr.FindAllByTemplateIDCall.Receives.Connection = conn
	kr.FindAllByTemplateIDCall.Receives.TemplateID = templateID

	return kr.FindAllByTemplateIDCall.Returns.Kinds, kr.FindAllByTemplateIDCall.Returns.Error
}

func (kr *KindsRepo) Trim(conn models.ConnectionInterface, clientID string, kindIDs []string) (int, error) {
	kr.TrimCall.Receives.Connection = conn
	kr.TrimCall.Receives.ClientID = clientID
	kr.TrimCall.Receives.KindIDs = kindIDs

	return kr.TrimCall.Returns.AffectedRowCount, kr.TrimCall.Returns.Error
}

func (kr *KindsRepo) Update(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	kr.UpdateCall.Receives.Connection = conn
	kr.UpdateCall.Receives.Kind = kind

	return kr.UpdateCall.Returns.Kind, kr.UpdateCall.Returns.Error
}

func (kr *KindsRepo) Upsert(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
	kr.UpsertCall.Receives.Connection = conn
	kr.UpsertCall.Receives.Kinds = append(kr.UpsertCall.Receives.Kinds, kind)

	return kr.UpsertCall.Returns.Kind, kr.UpsertCall.Returns.Error
}
