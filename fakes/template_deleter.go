package fakes

type FakeTemplateDeleter struct {
    DeleteArgument string
    DeleteError    error
}

func NewFakeTemplateDeleter() *FakeTemplateDeleter {
    return &FakeTemplateDeleter{}
}

func (fake *FakeTemplateDeleter) Delete(templateName string) error {
    fake.DeleteArgument = templateName
    return fake.DeleteError
}
