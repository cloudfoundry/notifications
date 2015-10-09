package v2

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
)

type NoStrategyError struct {
	Err error
}

func (e NoStrategyError) Error() string {
	return e.Err.Error()
}

type recipient struct {
	Email string
	GUID  string
}

type emailAddressFormatter interface {
	Format(email string) (formattedEmail string)
}

type htmlPartsExtractor interface {
	Extract(html string) (doctype, head, bodyContent, bodyAttributes string, err error)
}

type CampaignJobProcessor struct {
	emailFormatter emailAddressFormatter
	htmlExtractor  htmlPartsExtractor

	userStrategy  dispatcher
	spaceStrategy dispatcher
	orgStrategy   dispatcher
	emailStrategy dispatcher
}

type dispatcher interface {
	Dispatch(dispatch services.Dispatch) ([]services.Response, error)
}

func NewCampaignJobProcessor(emailFormatter emailAddressFormatter, htmlExtractor htmlPartsExtractor, userStrategy, spaceStrategy, orgStrategy, emailStrategy dispatcher) CampaignJobProcessor {
	return CampaignJobProcessor{
		emailFormatter: emailFormatter,
		htmlExtractor:  htmlExtractor,
		userStrategy:   userStrategy,
		spaceStrategy:  spaceStrategy,
		orgStrategy:    orgStrategy,
		emailStrategy:  emailStrategy,
	}
}

func (p CampaignJobProcessor) Process(conn services.ConnectionInterface, uaaHost string, job gobble.Job) error {
	var campaignJob queue.CampaignJob

	err := job.Unmarshal(&campaignJob)
	if err != nil {
		return err
	}

	var audience string
	for key, _ := range campaignJob.Campaign.SendTo {
		audience = key
	}

	doctype, head, bodyContent, bodyAttributes, err := p.htmlExtractor.Extract(campaignJob.Campaign.HTML)
	if err != nil {
		panic(err)
	}

	var recipients []recipient
	switch audience {
	case "emails":
		for _, audienceMember := range campaignJob.Campaign.SendTo[audience].([]interface{}) {
			recipients = append(recipients, recipient{
				Email: p.emailFormatter.Format(audienceMember.(string)),
			})
		}
	case "users", "spaces":
		for _, audienceMember := range campaignJob.Campaign.SendTo[audience].([]interface{}) {
			recipients = append(recipients, recipient{
				GUID: audienceMember.(string),
			})
		}
	default:
		recipients = append(recipients, recipient{
			GUID: campaignJob.Campaign.SendTo[audience].(string),
		})
	}

	strategy, err := p.findStrategy(audience)
	if err != nil {
		return err
	}

	for _, recipient := range recipients {
		_, err = strategy.Dispatch(services.Dispatch{
			JobType:    "v2",
			UAAHost:    uaaHost,
			GUID:       recipient.GUID,
			Connection: conn,
			TemplateID: campaignJob.Campaign.TemplateID,
			CampaignID: campaignJob.Campaign.ID,
			Client: services.DispatchClient{
				ID: campaignJob.Campaign.ClientID,
			},
			Message: services.DispatchMessage{
				To:      recipient.Email,
				ReplyTo: campaignJob.Campaign.ReplyTo,
				Subject: campaignJob.Campaign.Subject,
				Text:    campaignJob.Campaign.Text,
				HTML: services.HTML{
					Doctype:        doctype,
					Head:           head,
					BodyContent:    bodyContent,
					BodyAttributes: bodyAttributes,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p CampaignJobProcessor) findStrategy(audience string) (dispatcher, error) {
	switch audience {
	case "users":
		return p.userStrategy, nil
	case "spaces":
		return p.spaceStrategy, nil
	case "orgs":
		return p.orgStrategy, nil
	case "emails":
		return p.emailStrategy, nil
	default:
		return nil, NoStrategyError{fmt.Errorf("Strategy for the %q audience could not be found", audience)}
	}
}
