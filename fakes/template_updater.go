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

func (fake *TemplateUpdater) DeprecatedUpdate(template models.Template) error {
	fake.UpdateArgument = template

	return fake.UpdateError
}
