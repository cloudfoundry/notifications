package services

import "github.com/cloudfoundry-incubator/notifications/cf"

const SpaceEndorsement = `You received this message because you belong to the "{{.Space}}" space in the "{{.Organization}}" organization.`

type spaceUserGUIDFinder interface {
	UserGUIDsBelongingToSpace(spaceGUID, token string) (userGUIDs []string, err error)
}

type loadsSpaces interface {
	Load(spaceGUID, token string) (cf.CloudControllerSpace, error)
}

type SpaceStrategy struct {
	tokenLoader        loadsTokens
	spaceLoader        loadsSpaces
	organizationLoader loadsOrganizations
	findsUserGUIDs     spaceUserGUIDFinder
	v1Enqueuer         v1Enqueuer
	v2Enqueuer         v2Enqueuer
}

func NewSpaceStrategy(tokenLoader loadsTokens, spaceLoader loadsSpaces, organizationLoader loadsOrganizations, findsUserGUIDs spaceUserGUIDFinder, v1Enqueuer v1Enqueuer, v2Enqueuer v2Enqueuer) SpaceStrategy {
	return SpaceStrategy{
		tokenLoader:        tokenLoader,
		spaceLoader:        spaceLoader,
		organizationLoader: organizationLoader,
		findsUserGUIDs:     findsUserGUIDs,
		v1Enqueuer:         v1Enqueuer,
		v2Enqueuer:         v2Enqueuer,
	}
}

func (strategy SpaceStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	var responses []Response

	options := Options{
		To:                dispatch.Message.To,
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Endorsement:       SpaceEndorsement,
		Text:              dispatch.Message.Text,
		TemplateID:        dispatch.TemplateID,
		Role:              dispatch.Role,
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

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToSpace(dispatch.GUID, token)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	space, err := strategy.spaceLoader.Load(dispatch.GUID, token)
	if err != nil {
		return responses, err
	}

	org, err := strategy.organizationLoader.Load(space.OrganizationGUID, token)
	if err != nil {
		return responses, err
	}

	switch dispatch.JobType {
	case "v2":
		v2Users := convertToV2Users(users)
		v2Options := convertToV2Options(options)

		strategy.v2Enqueuer.Enqueue(dispatch.Connection, v2Users, v2Options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, dispatch.UAAHost, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)
	default:
		responses = strategy.v1Enqueuer.Enqueue(dispatch.Connection, users, options, space, org, dispatch.Client.ID, dispatch.UAAHost, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)
	}

	return responses, nil
}
