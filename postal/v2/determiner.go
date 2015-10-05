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

type emailAddressFormatter interface {
	Format(email string) (formattedEmail string)
}

type htmlPartsExtractor interface {
	Extract(html string) (doctype, head, bodyContent, bodyAttributes string, err error)
}

type Determiner struct {
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

func NewStrategyDeterminer(emailFormatter emailAddressFormatter, htmlExtractor htmlPartsExtractor, userStrategy, spaceStrategy, orgStrategy, emailStrategy dispatcher) Determiner {
	return Determiner{
		emailFormatter: emailFormatter,
		htmlExtractor:  htmlExtractor,
		userStrategy:   userStrategy,
		spaceStrategy:  spaceStrategy,
		orgStrategy:    orgStrategy,
		emailStrategy:  emailStrategy,
	}
}

func (d Determiner) Determine(conn services.ConnectionInterface, uaaHost string, job gobble.Job) error {
	var campaignJob queue.CampaignJob

	err := job.Unmarshal(&campaignJob)
	if err != nil {
		return err
	}

	var audience string
	for key, _ := range campaignJob.Campaign.SendTo {
		audience = key
	}

	doctype, head, bodyContent, bodyAttributes, err := d.htmlExtractor.Extract(campaignJob.Campaign.HTML)
	if err != nil {
		panic(err)
	}

	var recipient string
	if audience == "email" {
		recipient = d.emailFormatter.Format(campaignJob.Campaign.SendTo[audience])
	}

	strategy, err := d.findStrategy(audience)
	if err != nil {
		return err
	}

	var guid string
	if audience != "email" {
		guid = campaignJob.Campaign.SendTo[audience]
	}

	_, err = strategy.Dispatch(services.Dispatch{
		JobType:    "v2",
		UAAHost:    uaaHost,
		GUID:       guid,
		Connection: conn,
		TemplateID: campaignJob.Campaign.TemplateID,
		CampaignID: campaignJob.Campaign.ID,
		Client: services.DispatchClient{
			ID: campaignJob.Campaign.ClientID,
		},
		Message: services.DispatchMessage{
			To:      recipient,
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

	return nil
}

func (d Determiner) findStrategy(audience string) (dispatcher, error) {
	switch audience {
	case "user":
		return d.userStrategy, nil
	case "space":
		return d.spaceStrategy, nil
	case "org":
		return d.orgStrategy, nil
	case "email":
		return d.emailStrategy, nil
	default:
		return nil, NoStrategyError{fmt.Errorf("Strategy for the %q audience could not be found", audience)}
	}
}
