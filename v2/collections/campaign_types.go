package collections

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
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
	Insert(conn db.ConnectionInterface, campaignType models.CampaignType) (createdCampaignType models.CampaignType, err error)
	GetBySenderIDAndName(conn db.ConnectionInterface, senderID string, name string) (campaignType models.CampaignType, err error)
	List(conn db.ConnectionInterface, senderID string) (campaignTypes []models.CampaignType, err error)
	Get(conn db.ConnectionInterface, id string) (campaignType models.CampaignType, err error)
	Update(conn db.ConnectionInterface, campaignType models.CampaignType) (updatedCampaignType models.CampaignType, err error)
	Delete(conn db.ConnectionInterface, campaignType models.CampaignType) error
}

func NewCampaignTypesCollection(nr campaignTypesRepository, sr sendersRepository) CampaignTypesCollection {
	return CampaignTypesCollection{
		campaignTypesRepository: nr,
		sendersRepository:       sr,
	}
}

func (nc CampaignTypesCollection) Set(conn ConnectionInterface, campaignType CampaignType, clientID string) (CampaignType, error) {
	sender, err := nc.sendersRepository.Get(conn, campaignType.SenderID)
	err = validateSender(clientID, campaignType.SenderID, sender, err)
	if err != nil {
		return CampaignType{}, err
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
	}, nil
}

func (nc CampaignTypesCollection) Get(conn ConnectionInterface, campaignTypeID, senderID, clientID string) (CampaignType, error) {
	sender, err := nc.sendersRepository.Get(conn, senderID)
	err = validateSender(clientID, senderID, sender, err)
	if err != nil {
		return CampaignType{}, err
	}

	campaignType, err := nc.campaignTypesRepository.Get(conn, campaignTypeID)
	err = validateCampaignType(senderID, campaignTypeID, campaignType, err)
	if err != nil {
		return CampaignType{}, err
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

func (nc CampaignTypesCollection) List(conn ConnectionInterface, senderID, clientID string) ([]CampaignType, error) {
	sender, err := nc.sendersRepository.Get(conn, senderID)
	err = validateSender(clientID, senderID, sender, err)
	if err != nil {
		return []CampaignType{}, err
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

func (c CampaignTypesCollection) Delete(conn ConnectionInterface, campaignTypeID, senderID, clientID string) error {
	sender, err := c.sendersRepository.Get(conn, senderID)
	err = validateSender(clientID, senderID, sender, err)
	if err != nil {
		return err
	}

	campaignType, err := c.campaignTypesRepository.Get(conn, campaignTypeID)
	err = validateCampaignType(senderID, campaignTypeID, campaignType, err)
	if err != nil {
		return err
	}

	return c.campaignTypesRepository.Delete(conn, campaignType)
}

func validateSender(clientID, senderID string, sender models.Sender, err error) error {
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return NotFoundError{err}
		}

		return PersistenceError{err}
	}

	if sender.ClientID != clientID {
		return NotFoundError{errors.New(fmt.Sprintf("sender %s not found", senderID))}
	}

	return nil
}

func validateCampaignType(senderID, campaignTypeID string, campaignType models.CampaignType, err error) error {
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return NotFoundError{err}
		}

		return PersistenceError{err}
	}

	if campaignType.SenderID != senderID {
		return NotFoundError{errors.New(fmt.Sprintf("campaign type %s not found", campaignTypeID))}
	}

	return nil
}
