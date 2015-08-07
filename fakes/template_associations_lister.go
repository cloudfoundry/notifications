package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type TemplateAssociationLister struct {
	Associations map[string][]services.TemplateAssociation

	ListCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewTemplateAssociationLister() *TemplateAssociationLister {
	return &TemplateAssociationLister{
		Associations: make(map[string][]services.TemplateAssociation),
	}
}

func (lister *TemplateAssociationLister) List(database models.DatabaseInterface, templateID string) ([]services.TemplateAssociation, error) {
	lister.ListCall.Arguments = []interface{}{database, templateID}
	return lister.Associations[templateID], lister.ListCall.Error
}
