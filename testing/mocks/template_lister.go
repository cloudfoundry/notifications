package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/services"

type TemplateLister struct {
	ListCall struct {
		Receives struct {
			Database services.DatabaseInterface
		}
		Returns struct {
			TemplateSummaries map[string]services.TemplateSummary
			Error             error
		}
	}
}

func NewTemplateLister() *TemplateLister {
	return &TemplateLister{}
}

func (tl *TemplateLister) List(database services.DatabaseInterface) (map[string]services.TemplateSummary, error) {
	tl.ListCall.Receives.Database = database

	return tl.ListCall.Returns.TemplateSummaries, tl.ListCall.Returns.Error
}
