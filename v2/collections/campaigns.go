package collections

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type campaignEnqueuer interface {
	Enqueue(campaign Campaign, jobType string) error
}

type campaignsPersister interface {
	Insert(conn models.ConnectionInterface, campaign models.Campaign) (models.Campaign, error)
	Get(conn models.ConnectionInterface, campaignID string) (models.Campaign, error)
}

type campaignTypesGetter interface {
	Get(conn models.ConnectionInterface, campaignTypeID string) (models.CampaignType, error)
}

type templatesGetter interface {
	Get(conn models.ConnectionInterface, templateID string) (models.Template, error)
}

type sendersGetter interface {
	Get(conn models.ConnectionInterface, senderID string) (models.Sender, error)
}

type existenceChecker interface {
	Exists(guid string) (bool, error)
}

type Campaign struct {
	ID             string
	SendTo         map[string]string
	CampaignTypeID string
	Text           string
	HTML           string
	Subject        string
	TemplateID     string
	ReplyTo        string
	SenderID       string
	ClientID       string
}

type CampaignsCollection struct {
	enqueuer          campaignEnqueuer
	campaignsRepo     campaignsPersister
	campaignTypesRepo campaignTypesGetter
	templatesRepo     templatesGetter
	sendersRepo       sendersGetter
	userFinder        existenceChecker
	spaceFinder       existenceChecker
	orgFinder         existenceChecker
}

func NewCampaignsCollection(enqueuer campaignEnqueuer, campaignsRepo campaignsPersister, campaignTypesRepo campaignTypesGetter, templatesRepo templatesGetter, sendersRepo sendersGetter, userFinder existenceChecker, spaceFinder existenceChecker, orgFinder existenceChecker) CampaignsCollection {
	return CampaignsCollection{
		enqueuer:          enqueuer,
		campaignsRepo:     campaignsRepo,
		campaignTypesRepo: campaignTypesRepo,
		templatesRepo:     templatesRepo,
		sendersRepo:       sendersRepo,
		userFinder:        userFinder,
		spaceFinder:       spaceFinder,
		orgFinder:         orgFinder,
	}
}

func (c CampaignsCollection) Create(conn ConnectionInterface, campaign Campaign, clientID string, canSendCritical bool) (Campaign, error) {
	var audience string
	for key, _ := range campaign.SendTo {
		audience = key
	}

	exists, err := c.checkForExistence(audience, campaign.SendTo[audience])
	if err != nil {
		return Campaign{}, UnknownError{err}
	}

	if !exists {
		return Campaign{}, NotFoundError{fmt.Errorf("The %s %q cannot be found", audience, campaign.SendTo[audience])}
	}

	sender, err := c.sendersRepo.Get(conn, campaign.SenderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, UnknownError{err}
		}
	}

	if sender.ClientID != clientID {
		return Campaign{}, NotFoundError{fmt.Errorf("Sender with id %q could not be found", campaign.SenderID)}
	}

	campaignType, err := c.campaignTypesRepo.Get(conn, campaign.CampaignTypeID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, PersistenceError{err}
		}
	}

	if campaignType.Critical && !canSendCritical {
		return Campaign{}, PermissionsError{errors.New("Scope critical_notifications.write is required")}
	}

	if campaign.TemplateID == "" {
		campaign.TemplateID = campaignType.TemplateID
	}

	_, err = c.templatesRepo.Get(conn, campaign.TemplateID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, PersistenceError{err}
		}
	}

	sendTo, err := json.Marshal(campaign.SendTo)
	if err != nil {
		panic(err)
	}

	campaignModel, err := c.campaignsRepo.Insert(conn, models.Campaign{
		SendTo:         string(sendTo),
		CampaignTypeID: campaign.CampaignTypeID,
		Text:           campaign.Text,
		HTML:           campaign.HTML,
		Subject:        campaign.Subject,
		TemplateID:     campaign.TemplateID,
		ReplyTo:        campaign.ReplyTo,
		SenderID:       campaign.SenderID,
	})
	if err != nil {
		return Campaign{}, PersistenceError{err}
	}

	campaign.ID = campaignModel.ID
	campaign.ClientID = clientID

	err = c.enqueuer.Enqueue(campaign, "campaign")
	if err != nil {
		return Campaign{}, PersistenceError{Err: err}
	}

	return campaign, nil
}

func (c CampaignsCollection) checkForExistence(audience, guid string) (bool, error) {
	switch audience {
	case "user":
		return c.userFinder.Exists(guid)
	case "space":
		return c.spaceFinder.Exists(guid)
	case "org":
		return c.orgFinder.Exists(guid)
	case "email":
		return true, nil
	default:
		return false, fmt.Errorf("The %q audience is not valid", audience)
	}
}

func (c CampaignsCollection) Get(connection ConnectionInterface, campaignID, senderID, clientID string) (Campaign, error) {
	sender, err := c.sendersRepo.Get(connection, senderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, UnknownError{err}
		}
	}

	if sender.ClientID != clientID {
		return Campaign{}, NotFoundError{fmt.Errorf("Sender with id %q could not be found", senderID)}
	}

	model, err := c.campaignsRepo.Get(connection, campaignID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, UnknownError{err}
		}
	}

	if model.SenderID != senderID {
		return Campaign{}, NotFoundError{fmt.Errorf("Campaign with id %q could not be found", campaignID)}
	}

	var sendTo map[string]string
	err = json.Unmarshal([]byte(model.SendTo), &sendTo)
	if err != nil {
		panic(err)
	}

	return Campaign{
		ID:             campaignID,
		SendTo:         sendTo,
		CampaignTypeID: model.CampaignTypeID,
		Text:           model.Text,
		HTML:           model.HTML,
		Subject:        model.Subject,
		TemplateID:     model.TemplateID,
		ReplyTo:        model.ReplyTo,
		SenderID:       model.SenderID,
	}, nil
}
