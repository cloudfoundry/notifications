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
	Find(string) (models.Template, error)
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

func (finder TemplateFinder) Find(templateID string) (models.Template, error) {
	template, err := finder.search(finder.database.Connection(), templateID)
	if err != nil {
		return models.Template{}, err
	}

	return template, err
}

func (finder TemplateFinder) search(connection models.ConnectionInterface, templateID string) (models.Template, error) {
	template, err := finder.templatesRepo.Find(connection, templateID)

	if err != nil {
		if (err == models.ErrRecordNotFound{}) {
			return finder.findDefaultTemplate(templateID)
		}
		return models.Template{}, err
	}
	return template, nil
}

func (finder TemplateFinder) findDefaultTemplate(templateID string) (models.Template, error) {
	return finder.defaultTemplate(templateID)
}

func (finder TemplateFinder) defaultTemplate(name string) (models.Template, error) {
	text, err := finder.readTemplate(name + ".text")
	if err != nil {
		return models.Template{}, err
	}

	html, err := finder.readTemplate(name + ".html")
	if err != nil {
		return models.Template{}, err
	}

	return models.Template{Text: text, HTML: html}, nil
}

func (finder TemplateFinder) readTemplate(name string) (string, error) {
	path := finder.rootPath + "/templates/" + name
	return finder.fileSystem.Read(path)
}
