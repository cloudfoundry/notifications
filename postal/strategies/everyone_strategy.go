package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/services"
)

const EveryoneEndorsement = "This message was sent to everyone."

type EveryoneStrategy struct {
	tokenLoader TokenLoader
	allUsers    services.AllUsersInterface
	enqueuer    EnqueuerInterface
}

func NewEveryoneStrategy(tokenLoader TokenLoader, allUsers services.AllUsersInterface, enqueuer EnqueuerInterface) EveryoneStrategy {
	return EveryoneStrategy{
		tokenLoader: tokenLoader,
		allUsers:    allUsers,
		enqueuer:    enqueuer,
	}
}

func (strategy EveryoneStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	responses := []Response{}
	options := Options{
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		To:                dispatch.Message.To,
		Endorsement:       EveryoneEndorsement,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Text:              dispatch.Message.Text,
		HTML: HTML{
			BodyContent:    dispatch.Message.HTML.BodyContent,
			BodyAttributes: dispatch.Message.HTML.BodyAttributes,
			Head:           dispatch.Message.HTML.Head,
			Doctype:        dispatch.Message.HTML.Doctype,
		},
	}

	_, err := strategy.tokenLoader.Load()
	if err != nil {
		return responses, err
	}

	// split this up so that it only loads user guids
	userGUIDs, err := strategy.allUsers.AllUserGUIDs()
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	responses = strategy.enqueuer.Enqueue(dispatch.Connection, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)

	return responses, nil
}
