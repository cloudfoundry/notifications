package strategy

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
)

type Determiner struct {
	UserStrategy dispatcher
}

type dispatcher interface {
	Dispatch(dispatch services.Dispatch) ([]services.Response, error)
}

func (d Determiner) Determine(conn db.ConnectionInterface, uaaHost string, job gobble.Job) {
	var campaignJob queue.CampaignJob

	err := job.Unmarshal(&campaignJob)
	if err != nil {
		panic(err)
	}

	params := notify.NotifyParams{
		ReplyTo: campaignJob.Campaign.ReplyTo,
		Subject: campaignJob.Campaign.Subject,
		Text:    campaignJob.Campaign.Text,
		RawHTML: campaignJob.Campaign.HTML,
	}

	params.FormatEmailAndExtractHTML()

	_, err = d.UserStrategy.Dispatch(services.Dispatch{
		UAAHost:    uaaHost,
		GUID:       campaignJob.Campaign.SendTo["user"],
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
		panic(err)
	}
}
