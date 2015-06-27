package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/services"
)

const SpaceEndorsement = `You received this message because you belong to the "{{.Space}}" space in the "{{.Organization}}" organization.`

type SpaceStrategy struct {
	tokenLoader        postal.TokenLoaderInterface
	spaceLoader        services.SpaceLoaderInterface
	organizationLoader services.OrganizationLoaderInterface
	findsUserGUIDs     services.FindsUserGUIDsInterface
	mailer             MailerInterface
}

func NewSpaceStrategy(tokenLoader postal.TokenLoaderInterface, spaceLoader services.SpaceLoaderInterface, organizationLoader services.OrganizationLoaderInterface,
	findsUserGUIDs services.FindsUserGUIDsInterface, mailer MailerInterface) SpaceStrategy {

	return SpaceStrategy{
		tokenLoader:        tokenLoader,
		spaceLoader:        spaceLoader,
		organizationLoader: organizationLoader,
		findsUserGUIDs:     findsUserGUIDs,
		mailer:             mailer,
	}
}

func (strategy SpaceStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	responses := []Response{}
	options := postal.Options{
		To:                dispatch.Message.To,
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Endorsement:       SpaceEndorsement,
		Text:              dispatch.Message.Text,
		Role:              dispatch.Role,
		HTML: postal.HTML{
			BodyContent:    dispatch.Message.HTML.BodyContent,
			BodyAttributes: dispatch.Message.HTML.BodyAttributes,
			Head:           dispatch.Message.HTML.Head,
			Doctype:        dispatch.Message.HTML.Doctype,
		},
	}

	token, err := strategy.tokenLoader.Load()
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

	responses = strategy.mailer.Deliver(dispatch.Connection, users, options, space, org, dispatch.Client.ID, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)

	return responses, nil
}
