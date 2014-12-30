package fakes

type TemplateDeleter struct {
	DeleteArgument string
	DeleteError    error
}

func NewTemplateDeleter() *TemplateDeleter {
	return &TemplateDeleter{}
}

func (fake *TemplateDeleter) Delete(templateID string) error {
	fake.DeleteArgument = templateID
	return fake.DeleteError
}
