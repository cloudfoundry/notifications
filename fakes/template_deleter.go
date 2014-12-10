package fakes

type TemplateDeleter struct {
	DeleteArgument           string
	DeleteError              error
	DeprecatedDeleteArgument string
	DeprecatedDeleteError    error
}

func NewTemplateDeleter() *TemplateDeleter {
	return &TemplateDeleter{}
}

func (fake *TemplateDeleter) Delete(templateID string) error {
	fake.DeleteArgument = templateID
	return fake.DeleteError
}

func (fake *TemplateDeleter) DeprecatedDelete(templateName string) error {
	fake.DeprecatedDeleteArgument = templateName
	return fake.DeprecatedDeleteError
}
