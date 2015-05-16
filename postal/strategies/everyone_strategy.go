package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
)

const EveryoneEndorsement = "This message was sent to everyone."

type EveryoneStrategy struct {
	tokenLoader postal.TokenLoaderInterface
	allUsers    utilities.AllUsersInterface
	mailer      MailerInterface
}

func NewEveryoneStrategy(tokenLoader postal.TokenLoaderInterface, allUsers utilities.AllUsersInterface, mailer MailerInterface) EveryoneStrategy {
	return EveryoneStrategy{
		tokenLoader: tokenLoader,
		allUsers:    allUsers,
		mailer:      mailer,
	}
}

func (strategy EveryoneStrategy) Dispatch(clientID, guid, vcapRequestID string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	responses := []Response{}
	options.Endorsement = EveryoneEndorsement

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

	responses = strategy.mailer.Deliver(conn, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID, "", vcapRequestID)

	return responses, nil
}
