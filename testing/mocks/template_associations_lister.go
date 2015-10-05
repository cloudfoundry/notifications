package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/collections"

type TemplateAssociationLister struct {
	ListCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Associations []collections.TemplateAssociation
			Error        error
		}
	}
}

func NewTemplateAssociationLister() *TemplateAssociationLister {
	return &TemplateAssociationLister{}
}

func (l *TemplateAssociationLister) ListAssociations(connection collections.ConnectionInterface, templateID string) ([]collections.TemplateAssociation, error) {
	l.ListCall.Receives.Connection = connection
	l.ListCall.Receives.TemplateID = templateID

	return l.ListCall.Returns.Associations, l.ListCall.Returns.Error
}
