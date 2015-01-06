package fakes

import "github.com/cloudfoundry-incubator/notifications/web/services"

type TemplateAssociationLister struct {
	Associations map[string][]services.TemplateAssociation
	ListError    error
}

func NewTemplateAssociationLister() *TemplateAssociationLister {
	return &TemplateAssociationLister{
		Associations: make(map[string][]services.TemplateAssociation),
	}
}

func (lister *TemplateAssociationLister) List(templateID string) ([]services.TemplateAssociation, error) {
	return lister.Associations[templateID], lister.ListError
}
