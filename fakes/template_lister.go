package fakes

import "github.com/cloudfoundry-incubator/notifications/web/services"

type TemplateLister struct {
	ListWasCalled bool
	Templates     map[string]services.TemplateMetadata
	ListError     error
}

func NewTemplateLister() *TemplateLister {
	return &TemplateLister{}
}

func (lister *TemplateLister) List() (map[string]services.TemplateMetadata, error) {
	lister.ListWasCalled = true
	return lister.Templates, lister.ListError
}
