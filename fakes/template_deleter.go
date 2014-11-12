package fakes

type TemplateDeleter struct {
	DeleteArgument string
	DeleteError    error
}

func NewTemplateDeleter() *TemplateDeleter {
	return &TemplateDeleter{}
}

func (fake *TemplateDeleter) Delete(templateName string) error {
	fake.DeleteArgument = templateName
	return fake.DeleteError
}
