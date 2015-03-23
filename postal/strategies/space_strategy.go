package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
)

const SpaceEndorsement = `You received this message because you belong to the "{{.Space}}" space in the "{{.Organization}}" organization.`

type SpaceStrategy struct {
	tokenLoader        postal.TokenLoaderInterface
	spaceLoader        utilities.SpaceLoaderInterface
	organizationLoader utilities.OrganizationLoaderInterface
	findsUserGUIDs     utilities.FindsUserGUIDsInterface
	mailer             MailerInterface
}

func NewSpaceStrategy(tokenLoader postal.TokenLoaderInterface, spaceLoader utilities.SpaceLoaderInterface, organizationLoader utilities.OrganizationLoaderInterface,
	findsUserGUIDs utilities.FindsUserGUIDsInterface, mailer MailerInterface) SpaceStrategy {

	return SpaceStrategy{
		tokenLoader:        tokenLoader,
		spaceLoader:        spaceLoader,
		organizationLoader: organizationLoader,
		findsUserGUIDs:     findsUserGUIDs,
		mailer:             mailer,
	}
}

func (strategy SpaceStrategy) Dispatch(clientID, guid string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	responses := []Response{}
	options.Endorsement = SpaceEndorsement

	token, err := strategy.tokenLoader.Load()
	if err != nil {
		return responses, err
	}

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToSpace(guid, token)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	space, err := strategy.spaceLoader.Load(guid, token)
	if err != nil {
		return responses, err
	}

	org, err := strategy.organizationLoader.Load(space.OrganizationGUID, token)
	if err != nil {
		return responses, err
	}

	responses = strategy.mailer.Deliver(conn, users, options, space, org, clientID, "")

	return responses, nil
}
