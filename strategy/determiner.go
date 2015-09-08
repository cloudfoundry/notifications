package strategy

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
)

type NoStrategyError struct {
	Err error
}

func (e NoStrategyError) Error() string {
	return e.Err.Error()
}

type Determiner struct {
	userStrategy  dispatcher
	spaceStrategy dispatcher
	orgStrategy   dispatcher
	emailStrategy dispatcher
}

type dispatcher interface {
	Dispatch(dispatch services.Dispatch) ([]services.Response, error)
}

func NewStrategyDeterminer(userStrategy, spaceStrategy, orgStrategy, emailStrategy dispatcher) Determiner {
	return Determiner{
		userStrategy:  userStrategy,
		spaceStrategy: spaceStrategy,
		orgStrategy:   orgStrategy,
		emailStrategy: emailStrategy,
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

	params := notify.NotifyParams{
		ReplyTo: campaignJob.Campaign.ReplyTo,
		Subject: campaignJob.Campaign.Subject,
		Text:    campaignJob.Campaign.Text,
		RawHTML: campaignJob.Campaign.HTML,
	}
	if audience == "email" {
		params.To = campaignJob.Campaign.SendTo[audience]
	}

	params.FormatEmailAndExtractHTML()

	strategy, err := d.findStrategy(audience)
	if err != nil {
		return err
	}

	var guid string
	if audience != "email" {
		guid = campaignJob.Campaign.SendTo[audience]
	}

	_, err = strategy.Dispatch(services.Dispatch{
		UAAHost:    uaaHost,
		GUID:       guid,
		Connection: conn,
		TemplateID: campaignJob.Campaign.TemplateID,
		Client: services.DispatchClient{
			ID: campaignJob.Campaign.ClientID,
		},
		Message: services.DispatchMessage{
			To:      params.To,
			ReplyTo: params.ReplyTo,
			Subject: params.Subject,
			Text:    params.Text,
			HTML: services.HTML{
				BodyContent:    params.ParsedHTML.BodyContent,
				BodyAttributes: params.ParsedHTML.BodyAttributes,
				Head:           params.ParsedHTML.Head,
				Doctype:        params.ParsedHTML.Doctype,
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
