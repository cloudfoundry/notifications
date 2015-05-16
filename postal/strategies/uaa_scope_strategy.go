package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
)

const ScopeEndorsement = "You received this message because you have the {{.Scope}} scope."

type UAAScopeStrategy struct {
	findsUserGUIDs utilities.FindsUserGUIDsInterface
	tokenLoader    postal.TokenLoaderInterface
	mailer         MailerInterface
}

type DefaultScopeError struct{}

func (d DefaultScopeError) Error() string {
	return "You cannot send a notification to a default scope"
}

func NewUAAScopeStrategy(tokenLoader postal.TokenLoaderInterface, findsUserGUIDs utilities.FindsUserGUIDsInterface,
	mailer MailerInterface) UAAScopeStrategy {

	return UAAScopeStrategy{
		findsUserGUIDs: findsUserGUIDs,
		tokenLoader:    tokenLoader,
		mailer:         mailer,
	}
}

func (strategy UAAScopeStrategy) Dispatch(clientID, scope, vcapRequestID string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	responses := []Response{}
	options.Endorsement = ScopeEndorsement

	if strategy.scopeIsDefault(scope) {
		return responses, DefaultScopeError{}
	}

	_, err := strategy.tokenLoader.Load() // TODO: (rm) this triggers a weird side-effect that is required
	if err != nil {
		return responses, err
	}

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToScope(scope)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	responses = strategy.mailer.Deliver(conn, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID, scope, vcapRequestID)

	return responses, nil
}

func (strategy UAAScopeStrategy) scopeIsDefault(scope string) bool {
	defaultScopes := []string{
		"cloud_controller.read",
		"cloud_controller.write",
		"openid",
		"approvals.me",
		"cloud_controller_service_permissions.read",
		"scim.me",
		"uaa.user",
		"password.write",
		"scim.userids",
		"oauth.approvals",
	}

	for _, singleScope := range defaultScopes {
		if scope == singleScope {
			return true
		}
	}
	return false
}
