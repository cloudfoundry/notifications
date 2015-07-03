package services

import "github.com/cloudfoundry-incubator/notifications/cf"

const (
	OrganizationEndorsement     = `You received this message because you belong to the "{{.Organization}}" organization.`
	OrganizationRoleEndorsement = `You received this message because you are an {{.OrganizationRole}} in the "{{.Organization}}" organization.`
)

type OrganizationStrategy struct {
	tokenLoader        TokenLoader
	organizationLoader OrganizationLoaderInterface
	findsUserGUIDs     FindsUserGUIDsInterface
	enqueuer           EnqueuerInterface
}

func NewOrganizationStrategy(tokenLoader TokenLoader, organizationLoader OrganizationLoaderInterface, findsUserGUIDs FindsUserGUIDsInterface, enqueuer EnqueuerInterface) OrganizationStrategy {

	return OrganizationStrategy{
		tokenLoader:        tokenLoader,
		organizationLoader: organizationLoader,
		findsUserGUIDs:     findsUserGUIDs,
		enqueuer:           enqueuer,
	}
}

func (strategy OrganizationStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	responses := []Response{}
	options := Options{
		To:                dispatch.Message.To,
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Endorsement:       OrganizationEndorsement,
		Text:              dispatch.Message.Text,
		Role:              dispatch.Role,
		HTML: HTML{
			BodyContent:    dispatch.Message.HTML.BodyContent,
			BodyAttributes: dispatch.Message.HTML.BodyAttributes,
			Head:           dispatch.Message.HTML.Head,
			Doctype:        dispatch.Message.HTML.Doctype,
		},
	}

	if dispatch.Role != "" {
		options.Endorsement = OrganizationRoleEndorsement
	}

	token, err := strategy.tokenLoader.Load()
	if err != nil {
		return responses, err
	}

	organization, err := strategy.organizationLoader.Load(dispatch.GUID, token)
	if err != nil {
		return responses, err
	}

	userGUIDs, err := strategy.findsUserGUIDs.UserGUIDsBelongingToOrganization(dispatch.GUID, options.Role, token)
	if err != nil {
		return responses, err
	}

	var users []User
	for _, guid := range userGUIDs {
		users = append(users, User{GUID: guid})
	}

	responses = strategy.enqueuer.Enqueue(dispatch.Connection, users, options, cf.CloudControllerSpace{}, organization, dispatch.Client.ID, dispatch.UAAHost, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)

	return responses, nil
}
