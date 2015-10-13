package collections

import (
	"fmt"

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
	templatesRepository     templatesRepository
}

type campaignTypesRepository interface {
	Insert(conn models.ConnectionInterface, campaignType models.CampaignType) (createdCampaignType models.CampaignType, err error)
	GetBySenderIDAndName(conn models.ConnectionInterface, senderID string, name string) (campaignType models.CampaignType, err error)
	List(conn models.ConnectionInterface, senderID string) (campaignTypes []models.CampaignType, err error)
	Get(conn models.ConnectionInterface, id string) (campaignType models.CampaignType, err error)
	Update(conn models.ConnectionInterface, campaignType models.CampaignType) (updatedCampaignType models.CampaignType, err error)
	Delete(conn models.ConnectionInterface, campaignType models.CampaignType) error
}

func NewCampaignTypesCollection(nr campaignTypesRepository, sr sendersRepository, tr templatesRepository) CampaignTypesCollection {
	return CampaignTypesCollection{
		campaignTypesRepository: nr,
		sendersRepository:       sr,
		templatesRepository:     tr,
	}
}

func (nc CampaignTypesCollection) Set(conn ConnectionInterface, campaignType CampaignType, clientID string) (CampaignType, error) {
	sender, err := nc.sendersRepository.Get(conn, campaignType.SenderID)
	err = validateSender(clientID, campaignType.SenderID, sender, err)
	if err != nil {
		return CampaignType{}, err
	}

	if campaignType.TemplateID != "" {
		template, err := nc.templatesRepository.Get(conn, campaignType.TemplateID)
		err = validateTemplate(clientID, template, err)
		if err != nil {
			return CampaignType{}, err
		}
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

func (nc CampaignTypesCollection) Get(conn ConnectionInterface, campaignTypeID, clientID string) (CampaignType, error) {
	campaignType, err := nc.campaignTypesRepository.Get(conn, campaignTypeID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return CampaignType{}, NotFoundError{err}
		}

		return CampaignType{}, PersistenceError{err}
	}

	sender, err := nc.sendersRepository.Get(conn, campaignType.SenderID)
	err = validateSender(clientID, campaignType.SenderID, sender, err)
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

func (c CampaignTypesCollection) Delete(conn ConnectionInterface, campaignTypeID, clientID string) error {
	campaignType, err := c.campaignTypesRepository.Get(conn, campaignTypeID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return NotFoundError{err}
		}

		return PersistenceError{err}
	}

	sender, err := c.sendersRepository.Get(conn, campaignType.SenderID)
	err = validateSender(clientID, campaignType.SenderID, sender, err)
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
		return NotFoundError{fmt.Errorf("Sender with id %q could not be found", senderID)}
	}

	return nil
}

func validateTemplate(clientID string, template models.Template, err error) error {
	if err != nil {
		if _, ok := err.(models.RecordNotFoundError); ok {
			return NotFoundError{err}
		}

		return PersistenceError{err}
	}

	if clientID != template.ClientID {
		return NotFoundError{fmt.Errorf("Template with id %q could not be found", template.ID)}
	}

	return nil
}
