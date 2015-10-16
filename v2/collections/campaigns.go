package collections

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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

type Campaign struct {
	ID             string
	SendTo         map[string][]string
	CampaignTypeID string
	Text           string
	HTML           string
	Subject        string
	TemplateID     string
	ReplyTo        string
	SenderID       string
	ClientID       string
	StartTime      time.Time
}

type CampaignsCollection struct {
	enqueuer          campaignEnqueuer
	campaignsRepo     campaignsPersister
	campaignTypesRepo campaignTypesGetter
	templatesRepo     templatesGetter
	sendersRepo       sendersGetter
}

func NewCampaignsCollection(enqueuer campaignEnqueuer, campaignsRepo campaignsPersister, campaignTypesRepo campaignTypesGetter, templatesRepo templatesGetter, sendersRepo sendersGetter) CampaignsCollection {
	return CampaignsCollection{
		enqueuer:          enqueuer,
		campaignsRepo:     campaignsRepo,
		campaignTypesRepo: campaignTypesRepo,
		templatesRepo:     templatesRepo,
		sendersRepo:       sendersRepo,
	}
}

func (c CampaignsCollection) Create(conn ConnectionInterface, campaign Campaign, clientID string, canSendCritical bool) (Campaign, error) {
	for audience, audienceMembers := range campaign.SendTo {
		for _, audienceMember := range audienceMembers {
			exists, err := c.checkForExistence(audience, audienceMember)
			if err != nil {
				return Campaign{}, UnknownError{err}
			}

			if !exists {
				return Campaign{}, NotFoundError{fmt.Errorf("The %s %q cannot be found", strings.TrimSuffix(audience, "s"), audienceMember)}
			}
		}
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

	if campaign.TemplateID == "" {
		campaign.TemplateID = models.DefaultTemplate.ID
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
		StartTime:      campaign.StartTime,
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
	case "users":
		return true, nil
	case "spaces":
		return true, nil
	case "orgs":
		return true, nil
	case "emails":
		return true, nil
	default:
		return false, fmt.Errorf("The %q audience is not valid", audience)
	}
}

func (c CampaignsCollection) Get(connection ConnectionInterface, campaignID, clientID string) (Campaign, error) {
	campaign, err := c.campaignsRepo.Get(connection, campaignID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, UnknownError{err}
		}
	}

	sender, err := c.sendersRepo.Get(connection, campaign.SenderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, UnknownError{err}
		}
	}

	if sender.ClientID != clientID {
		return Campaign{}, NotFoundError{fmt.Errorf("Campaign with id %q could not be found", campaignID)}
	}

	var sendTo map[string][]string
	err = json.Unmarshal([]byte(campaign.SendTo), &sendTo)
	if err != nil {
		panic(err)
	}

	return Campaign{
		ID:             campaignID,
		SendTo:         sendTo,
		CampaignTypeID: campaign.CampaignTypeID,
		Text:           campaign.Text,
		HTML:           campaign.HTML,
		Subject:        campaign.Subject,
		TemplateID:     campaign.TemplateID,
		ReplyTo:        campaign.ReplyTo,
		SenderID:       campaign.SenderID,
	}, nil
}
