package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdater struct {
	UpdateArgumentID   string
	UpdateArgumentBody models.Template
	UpdateArgument     models.Template
	UpdateError        error
}

func NewTemplateUpdater() *TemplateUpdater {
	return &TemplateUpdater{}
}

func (fake *TemplateUpdater) Update(templateID string, template models.Template) error {
	fake.UpdateArgumentID = templateID
	fake.UpdateArgumentBody = template

	return fake.UpdateError
}
