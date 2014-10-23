package models

type TemplatesRepoInterface interface {
    Find(string) (Template, error)
}

type TemplatesRepo struct{}

func NewTemplatesRepo() TemplatesRepo {
    return TemplatesRepo{}
}

func (repo TemplatesRepo) Find(templateName string) (Template, error) {
    return Template{}, ErrRecordNotFound{}
}
