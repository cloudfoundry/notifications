package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
)

const (
	OrganizationEndorsement     = `You received this message because you belong to the "{{.Organization}}" organization.`
	OrganizationRoleEndorsement = `You received this message because you are an {{.OrganizationRole}} in the "{{.Organization}}" organization.`
)

type OrganizationStrategy struct {
	tokenLoader        postal.TokenLoaderInterface
	organizationLoader utilities.OrganizationLoaderInterface
	findsUserGUIDs     utilities.FindsUserGUIDsInterface
	mailer             MailerInterface
}

func NewOrganizationStrategy(tokenLoader postal.TokenLoaderInterface, organizationLoader utilities.OrganizationLoaderInterface,
	findsUserGUIDs utilities.FindsUserGUIDsInterface, mailer MailerInterface) OrganizationStrategy {

	return OrganizationStrategy{
		tokenLoader:        tokenLoader,
		organizationLoader: organizationLoader,
		findsUserGUIDs:     findsUserGUIDs,
		mailer:             mailer,
	}
}

func (strategy OrganizationStrategy) Dispatch(clientID, guid string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	responses := []Response{}
	options.Endorsement = OrganizationEndorsement

	if options.Role != "" {
		options.Endorsement = OrganizationRoleEndorsement
	}

	token, err := strategy.tokenLoader.Load()
	if err != nil {
		return responses, err
	}

	organization, err := strategy.organizationLoader.Load(guid, token)
	if err != nil {
		return responses, err
	}

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToOrganization(guid, options.Role, token)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	responses = strategy.mailer.Deliver(conn, users, options, cf.CloudControllerSpace{}, organization, clientID, "")

	return responses, nil
}
