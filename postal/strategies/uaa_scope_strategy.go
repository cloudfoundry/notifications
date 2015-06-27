package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/services"
)

const ScopeEndorsement = "You received this message because you have the {{.Scope}} scope."

type UAAScopeStrategy struct {
	findsUserGUIDs services.FindsUserGUIDsInterface
	tokenLoader    postal.TokenLoaderInterface
	mailer         MailerInterface
}

type DefaultScopeError struct{}

func (d DefaultScopeError) Error() string {
	return "You cannot send a notification to a default scope"
}

func NewUAAScopeStrategy(tokenLoader postal.TokenLoaderInterface, findsUserGUIDs services.FindsUserGUIDsInterface,
	mailer MailerInterface) UAAScopeStrategy {

	return UAAScopeStrategy{
		findsUserGUIDs: findsUserGUIDs,
		tokenLoader:    tokenLoader,
		mailer:         mailer,
	}
}

func (strategy UAAScopeStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	responses := []Response{}
	options := postal.Options{
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		To:                dispatch.Message.To,
		Endorsement:       ScopeEndorsement,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Text:              dispatch.Message.Text,
		HTML: postal.HTML{
			BodyContent:    dispatch.Message.HTML.BodyContent,
			BodyAttributes: dispatch.Message.HTML.BodyAttributes,
			Head:           dispatch.Message.HTML.Head,
			Doctype:        dispatch.Message.HTML.Doctype,
		},
	}

	if strategy.scopeIsDefault(dispatch.GUID) {
		return responses, DefaultScopeError{}
	}

	_, err := strategy.tokenLoader.Load() // TODO: (rm) this triggers a weird side-effect that is required
	if err != nil {
		return responses, err
	}

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToScope(dispatch.GUID)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	responses = strategy.mailer.Deliver(dispatch.Connection, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, dispatch.GUID, dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)

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
