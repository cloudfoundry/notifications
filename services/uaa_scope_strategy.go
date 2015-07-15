package services

import "github.com/cloudfoundry-incubator/notifications/cf"

const ScopeEndorsement = "You received this message because you have the {{.Scope}} scope."

type UAAScopeStrategy struct {
	findsUserGUIDs FindsUserGUIDsInterface
	tokenLoader    ZonedTokenLoaderInterface
	enqueuer       EnqueuerInterface
	defaultScopes  []string
}

type DefaultScopeError struct{}

func (d DefaultScopeError) Error() string {
	return "You cannot send a notification to a default scope"
}

func NewUAAScopeStrategy(tokenLoader ZonedTokenLoaderInterface, findsUserGUIDs FindsUserGUIDsInterface, enqueuer EnqueuerInterface, defaultScopes []string) UAAScopeStrategy {
	return UAAScopeStrategy{
		findsUserGUIDs: findsUserGUIDs,
		tokenLoader:    tokenLoader,
		enqueuer:       enqueuer,
		defaultScopes:  defaultScopes,
	}
}

func (strategy UAAScopeStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	responses := []Response{}
	options := Options{
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		To:                dispatch.Message.To,
		Endorsement:       ScopeEndorsement,
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

	if strategy.scopeIsDefault(dispatch.GUID) {
		return responses, DefaultScopeError{}
	}

	token, err := strategy.tokenLoader.Load(dispatch.UAAHost) // TODO: (rm) this triggers a weird side-effect that is required
	if err != nil {
		return responses, err
	}

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToScope(token, dispatch.GUID)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	responses = strategy.enqueuer.Enqueue(dispatch.Connection, users, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, dispatch.Client.ID, dispatch.UAAHost, dispatch.GUID, dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)

	return responses, nil
}

func (strategy UAAScopeStrategy) scopeIsDefault(scope string) bool {
	for _, singleScope := range strategy.defaultScopes {
		if scope == singleScope {
			return true
		}
	}
	return false
}
