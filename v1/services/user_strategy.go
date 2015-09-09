package services

import "github.com/cloudfoundry-incubator/notifications/cf"

const UserEndorsement = "This message was sent directly to you."

type UserStrategy struct {
	v1Enqueuer v1Enqueuer
	v2Enqueuer v2Enqueuer
}

func NewUserStrategy(v1Enqueuer v1Enqueuer, v2Enqueuer v2Enqueuer) UserStrategy {
	return UserStrategy{
		v1Enqueuer: v1Enqueuer,
		v2Enqueuer: v2Enqueuer,
	}
}

func (strategy UserStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	var responses []Response

	options := Options{
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		To:                dispatch.Message.To,
		Endorsement:       UserEndorsement,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Text:              dispatch.Message.Text,
		TemplateID:        dispatch.TemplateID,
		HTML: HTML{
			BodyContent:    dispatch.Message.HTML.BodyContent,
			BodyAttributes: dispatch.Message.HTML.BodyAttributes,
			Head:           dispatch.Message.HTML.Head,
			Doctype:        dispatch.Message.HTML.Doctype,
		},
	}
	users := []User{{GUID: dispatch.GUID}}

	switch dispatch.JobType {
	case "v2":
		v2Users := convertToV2Users(users)
		v2Options := convertToV2Options(options)

		strategy.v2Enqueuer.Enqueue(dispatch.Connection, v2Users, v2Options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, dispatch.UAAHost, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)
	default:
		responses = strategy.v1Enqueuer.Enqueue(dispatch.Connection, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, dispatch.UAAHost, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)
	}

	return responses, nil
}
