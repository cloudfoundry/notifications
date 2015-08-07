package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type TemplateLister struct {
	Templates map[string]services.TemplateSummary

	ListCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewTemplateLister() *TemplateLister {
	return &TemplateLister{}
}

func (lister *TemplateLister) List(database models.DatabaseInterface) (map[string]services.TemplateSummary, error) {
	lister.ListCall.Arguments = []interface{}{database}
	return lister.Templates, lister.ListCall.Error
}
