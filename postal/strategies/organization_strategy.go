package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/services"
)

const (
	OrganizationEndorsement     = `You received this message because you belong to the "{{.Organization}}" organization.`
	OrganizationRoleEndorsement = `You received this message because you are an {{.OrganizationRole}} in the "{{.Organization}}" organization.`
)

type OrganizationStrategy struct {
	tokenLoader        postal.TokenLoaderInterface
	organizationLoader services.OrganizationLoaderInterface
	findsUserGUIDs     services.FindsUserGUIDsInterface
	mailer             MailerInterface
}

func NewOrganizationStrategy(tokenLoader postal.TokenLoaderInterface, organizationLoader services.OrganizationLoaderInterface,
	findsUserGUIDs services.FindsUserGUIDsInterface, mailer MailerInterface) OrganizationStrategy {

	return OrganizationStrategy{
		tokenLoader:        tokenLoader,
		organizationLoader: organizationLoader,
		findsUserGUIDs:     findsUserGUIDs,
		mailer:             mailer,
	}
}

func (strategy OrganizationStrategy) Dispatch(dispatch Dispatch) ([]Response, error) {
	responses := []Response{}
	options := postal.Options{
		To:                dispatch.Message.To,
		ReplyTo:           dispatch.Message.ReplyTo,
		Subject:           dispatch.Message.Subject,
		KindID:            dispatch.Kind.ID,
		KindDescription:   dispatch.Kind.Description,
		SourceDescription: dispatch.Client.Description,
		Endorsement:       OrganizationEndorsement,
		Text:              dispatch.Message.Text,
		Role:              dispatch.Role,
		HTML: postal.HTML{
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

	responses = strategy.mailer.Deliver(dispatch.Connection, users, options, cf.CloudControllerSpace{}, organization, dispatch.Client.ID, "", dispatch.VCAPRequest.ID, dispatch.VCAPRequest.ReceiptTime)

	return responses, nil
}
