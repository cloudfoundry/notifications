package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

const EmailEndorsement = "This message was sent directly to your email address."

type EmailStrategy struct {
	mailer MailerInterface
}

func NewEmailStrategy(mailer MailerInterface) EmailStrategy {
	return EmailStrategy{
		mailer: mailer,
	}
}

func (strategy EmailStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	options := postal.Options{
		To:                dispatch.Message.To,
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Endorsement:       EmailEndorsement,
		Text:              dispatch.Message.Text,
		HTML: postal.HTML{
			BodyContent:    dispatch.Message.HTML.BodyContent,
			BodyAttributes: dispatch.Message.HTML.BodyAttributes,
			Head:           dispatch.Message.HTML.Head,
			Doctype:        dispatch.Message.HTML.Doctype,
		},
	}
	users := []User{{Email: dispatch.Message.To}}
	responses := strategy.mailer.Deliver(dispatch.Connection, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)

	return responses, nil
}
