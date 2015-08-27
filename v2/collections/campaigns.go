package collections

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type campaignEnqueuer interface {
	Enqueue(campaign Campaign, jobType string) error
}

type campaignTypesGetter interface {
	Get(conn models.ConnectionInterface, campaignTypeID string) (models.CampaignType, error)
}

type templatesGetter interface {
	Get(conn models.ConnectionInterface, templateID string) (models.Template, error)
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
	ClientID       string
}

type CampaignsCollection struct {
	enqueuer          campaignEnqueuer
	campaignTypesRepo campaignTypesGetter
	templatesRepo     templatesGetter
	userFinder        existenceChecker
	spaceFinder       existenceChecker
	orgFinder         existenceChecker
}

func NewCampaignsCollection(enqueuer campaignEnqueuer, campaignTypesRepo campaignTypesGetter, templatesRepo templatesGetter, userFinder existenceChecker, spaceFinder existenceChecker, orgFinder existenceChecker) CampaignsCollection {
	return CampaignsCollection{
		enqueuer:          enqueuer,
		campaignTypesRepo: campaignTypesRepo,
		templatesRepo:     templatesRepo,
		userFinder:        userFinder,
		spaceFinder:       spaceFinder,
		orgFinder:         orgFinder,
	}
}

func (c CampaignsCollection) Create(conn ConnectionInterface, campaign Campaign, canSendCritical bool) (Campaign, error) {
	campaign.ID = "some-random-id"
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
	default:
		return false, fmt.Errorf("The %q audience is not valid", audience)
	}
}
