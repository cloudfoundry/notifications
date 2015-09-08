package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/services"

type TemplateAssociationLister struct {
	ListCall struct {
		Receives struct {
			Database   services.DatabaseInterface
			TemplateID string
		}
		Returns struct {
			Associations []services.TemplateAssociation
			Error        error
		}
	}
}

func NewTemplateAssociationLister() *TemplateAssociationLister {
	return &TemplateAssociationLister{}
}

func (l *TemplateAssociationLister) List(database services.DatabaseInterface, templateID string) ([]services.TemplateAssociation, error) {
	l.ListCall.Receives.Database = database
	l.ListCall.Receives.TemplateID = templateID

	return l.ListCall.Returns.Associations, l.ListCall.Returns.Error
}
