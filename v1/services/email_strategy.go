package services

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
)

const EmailEndorsement = "This message was sent directly to your email address."

type EmailStrategy struct {
	v1Enqueuer v1Enqueuer
	v2Enqueuer v2Enqueuer
}

type v1Enqueuer interface {
	Enqueue(conn ConnectionInterface, users []User, opts Options, space cf.CloudControllerSpace, org cf.CloudControllerOrganization, clientID, uaaHost, scope, vcapRequestID string, reqReceived time.Time) []Response
}

type v2Enqueuer interface {
	Enqueue(conn queue.ConnectionInterface, users []queue.User, opts queue.Options, space cf.CloudControllerSpace, org cf.CloudControllerOrganization, clientID, uaaHost, scope, vcapRequestID string, reqReceived time.Time, campaignID string) []queue.Response
}

func NewEmailStrategy(v1Enqueuer v1Enqueuer, v2Enqueuer v2Enqueuer) EmailStrategy {
	return EmailStrategy{
		v1Enqueuer: v1Enqueuer,
		v2Enqueuer: v2Enqueuer,
	}
}

func (strategy EmailStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	var responses []Response

	options := Options{
		To:                dispatch.Message.To,
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Endorsement:       EmailEndorsement,
		Text:              dispatch.Message.Text,
		TemplateID:        dispatch.TemplateID,
		HTML: HTML{
			BodyContent:    dispatch.Message.HTML.BodyContent,
			BodyAttributes: dispatch.Message.HTML.BodyAttributes,
			Head:           dispatch.Message.HTML.Head,
			Doctype:        dispatch.Message.HTML.Doctype,
		},
	}
	users := []User{{Email: dispatch.Message.To}}

	switch dispatch.JobType {
	case "v2":
		v2Users := convertToV2Users(users)
		v2Options := convertToV2Options(options)

		strategy.v2Enqueuer.Enqueue(dispatch.Connection, v2Users, v2Options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, dispatch.UAAHost, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime, dispatch.CampaignID)
	default:
		responses = strategy.v1Enqueuer.Enqueue(dispatch.Connection, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, dispatch.UAAHost, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)
	}

	return responses, nil
}

func convertToV2Users(users []User) []queue.User {
	var v2Users []queue.User
	for _, user := range users {
		v2Users = append(v2Users, queue.User{
			GUID:  user.GUID,
			Email: user.Email,
		})
	}
	return v2Users
}

func convertToV2Options(options Options) queue.Options {
	return queue.Options{
		ReplyTo:           options.ReplyTo,
		Subject:           options.Subject,
		KindDescription:   options.KindDescription,
		SourceDescription: options.SourceDescription,
		Text:              options.Text,
		HTML: queue.HTML{
			BodyContent:    options.HTML.BodyContent,
			BodyAttributes: options.HTML.BodyAttributes,
			Head:           options.HTML.Head,
			Doctype:        options.HTML.Doctype,
		},
		KindID:      options.KindID,
		To:          options.To,
		Role:        options.Role,
		Endorsement: options.Endorsement,
		TemplateID:  options.TemplateID,
	}
}
