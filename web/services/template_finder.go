package services

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type FileSystemInterface interface {
	Exists(string) bool
	Read(string) (string, error)
}

type TemplateFinder struct {
	templatesRepo models.TemplatesRepoInterface
	rootPath      string
	fileSystem    FileSystemInterface
	database      models.DatabaseInterface
}

type TemplateFinderInterface interface {
	FindByID(string) (models.Template, error)
}

type TemplateNotFoundError string

func (err TemplateNotFoundError) Error() string {
	return fmt.Sprintf("Template %q could not be found.", string(err))
}

func NewTemplateFinder(templatesRepo models.TemplatesRepoInterface, rootPath string, database models.DatabaseInterface, fileSystem FileSystemInterface) TemplateFinder {
	return TemplateFinder{
		templatesRepo: templatesRepo,
		rootPath:      rootPath,
		database:      database,
		fileSystem:    fileSystem,
	}
}

func (finder TemplateFinder) FindByID(templateID string) (models.Template, error) {
	template, err := finder.templatesRepo.FindByID(finder.database.Connection(), templateID)
	if err != nil {
		return models.Template{}, err
	}

	return template, err
}
