package collections

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type CampaignType struct {
	ID          string
	Name        string
	Description string
	Critical    bool
	TemplateID  string
	SenderID    string
}

type CampaignTypesCollection struct {
	campaignTypesRepository campaignTypesRepository
	sendersRepository       sendersRepository
}

type campaignTypesRepository interface {
	Insert(conn models.ConnectionInterface, campaignType models.CampaignType) (createdCampaignType models.CampaignType, err error)
	GetBySenderIDAndName(conn models.ConnectionInterface, senderID string, name string) (campaignType models.CampaignType, err error)
	List(conn models.ConnectionInterface, senderID string) (campaignTypes []models.CampaignType, err error)
	Get(conn models.ConnectionInterface, id string) (campaignType models.CampaignType, err error)
	Update(conn models.ConnectionInterface, campaignType models.CampaignType) (updatedCampaignType models.CampaignType, err error)
}

func NewCampaignTypesCollection(nr campaignTypesRepository, sr sendersRepository) CampaignTypesCollection {
	return CampaignTypesCollection{
		campaignTypesRepository: nr,
		sendersRepository:       sr,
	}
}

func (nc CampaignTypesCollection) Set(conn models.ConnectionInterface, campaignType CampaignType, clientID string) (CampaignType, error) {
	senderModel, err := nc.sendersRepository.Get(conn, campaignType.SenderID)
	if err != nil {
		switch e := err.(type) {
		case models.RecordNotFoundError:
			return CampaignType{}, NotFoundError{
				Err:     e,
				Message: fmt.Sprintf("Sender %s not found", campaignType.SenderID),
			}
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	if senderModel.ClientID != clientID {
		return CampaignType{}, NewNotFoundError(fmt.Sprintf("Sender %s not found", campaignType.SenderID))
	}

	var (
		returnCampaignType models.CampaignType
		campaignTypeModel  = models.CampaignType{
			ID:          campaignType.ID,
			Name:        campaignType.Name,
			Description: campaignType.Description,
			Critical:    campaignType.Critical,
			TemplateID:  campaignType.TemplateID,
			SenderID:    campaignType.SenderID,
		}
	)

	if campaignType.ID != "" {
		returnCampaignType, err = nc.campaignTypesRepository.Update(conn, campaignTypeModel)
	} else {
		returnCampaignType, err = nc.campaignTypesRepository.Insert(conn, campaignTypeModel)
	}
	if err != nil {
		switch err.(type) {
		case models.DuplicateRecordError:
			returnCampaignType, err = nc.campaignTypesRepository.GetBySenderIDAndName(conn, campaignType.SenderID, campaignType.Name)
			if err != nil {
				return CampaignType{}, PersistenceError{err}
			}
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	return CampaignType{
		ID:          returnCampaignType.ID,
		Name:        returnCampaignType.Name,
		Description: returnCampaignType.Description,
		Critical:    returnCampaignType.Critical,
		TemplateID:  returnCampaignType.TemplateID,
		SenderID:    returnCampaignType.SenderID,
	}, err
}

func (nc CampaignTypesCollection) Get(conn models.ConnectionInterface, campaignTypeID, senderID, clientID string) (CampaignType, error) {
	sender, err := nc.sendersRepository.Get(conn, senderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return CampaignType{}, NewNotFoundError(fmt.Sprintf("sender %s not found", campaignTypeID))
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	if clientID != sender.ClientID {
		return CampaignType{}, NewNotFoundError(fmt.Sprintf("sender %s not found", senderID))
	}

	campaignType, err := nc.campaignTypesRepository.Get(conn, campaignTypeID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return CampaignType{}, NotFoundError{
				Err:     err,
				Message: fmt.Sprintf("campaign type %s not found", campaignTypeID),
			}
		default:
			return CampaignType{}, PersistenceError{err}
		}
	}

	if senderID != campaignType.SenderID {
		return CampaignType{}, NewNotFoundError(fmt.Sprintf("campaign type %s not found", campaignTypeID))
	}

	return CampaignType{
		ID:          campaignType.ID,
		Name:        campaignType.Name,
		Description: campaignType.Description,
		Critical:    campaignType.Critical,
		TemplateID:  campaignType.TemplateID,
		SenderID:    campaignType.SenderID,
	}, nil
}

func (nc CampaignTypesCollection) List(conn models.ConnectionInterface, senderID, clientID string) ([]CampaignType, error) {
	senderModel, err := nc.sendersRepository.Get(conn, senderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return []CampaignType{}, NotFoundError{
				Err:     err,
				Message: "sender not found",
			}
		default:
			return []CampaignType{}, PersistenceError{err}
		}
	}

	if senderModel.ClientID != clientID {
		return []CampaignType{}, NewNotFoundError("sender not found")
	}

	modelList, err := nc.campaignTypesRepository.List(conn, senderID)
	if err != nil {
		return []CampaignType{}, PersistenceError{err}
	}

	campaignTypeList := []CampaignType{}

	for _, model := range modelList {
		campaignType := CampaignType{
			ID:          model.ID,
			Name:        model.Name,
			Description: model.Description,
			Critical:    model.Critical,
			TemplateID:  model.TemplateID,
			SenderID:    model.SenderID,
		}
		campaignTypeList = append(campaignTypeList, campaignType)
	}

	return campaignTypeList, nil
}
