package fakes

import "github.com/cloudfoundry-incubator/notifications/web/services"

type TemplateLister struct {
	ListWasCalled bool
	Templates     map[string]services.TemplateSummary
	ListError     error
}

func NewTemplateLister() *TemplateLister {
	return &TemplateLister{}
}

func (lister *TemplateLister) List() (map[string]services.TemplateSummary, error) {
	lister.ListWasCalled = true
	return lister.Templates, lister.ListError
}
