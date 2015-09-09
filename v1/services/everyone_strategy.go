package services

import "github.com/cloudfoundry-incubator/notifications/cf"

const EveryoneEndorsement = "This message was sent to everyone."

type allUserGUIDsGetter interface {
	AllUserGUIDs(token string) (userGUIDs []string, err error)
}

type loadsTokens interface {
	Load(host string) (token string, err error)
}

type EveryoneStrategy struct {
	tokenLoader loadsTokens
	allUsers    allUserGUIDsGetter
	v1Enqueuer  v1Enqueuer
	v2Enqueuer  v2Enqueuer
}

func NewEveryoneStrategy(tokenLoader loadsTokens, allUsers allUserGUIDsGetter, v1Enqueuer v1Enqueuer, v2Enqueuer v2Enqueuer) EveryoneStrategy {
	return EveryoneStrategy{
		tokenLoader: tokenLoader,
		allUsers:    allUsers,
		v1Enqueuer:  v1Enqueuer,
		v2Enqueuer:  v2Enqueuer,
	}
}

func (strategy EveryoneStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	var responses []Response

	options := Options{
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		To:                dispatch.Message.To,
		Endorsement:       EveryoneEndorsement,
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

	token, err := strategy.tokenLoader.Load(dispatch.UAAHost)
	if err != nil {
		return responses, err
	}

	// split this up so that it only loads user guids
	userGUIDs, err := strategy.allUsers.AllUserGUIDs(token)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

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
