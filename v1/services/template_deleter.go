package services

type TemplateDeleter struct {
	templatesRepo TemplatesRepo
}

func NewTemplateDeleter(templatesRepo TemplatesRepo) TemplateDeleter {
	return TemplateDeleter{
		templatesRepo: templatesRepo,
	}
}

func (deleter TemplateDeleter) Delete(database DatabaseInterface, templateID string) error {
	connection := database.Connection()
	return deleter.templatesRepo.Destroy(connection, templateID)
}
