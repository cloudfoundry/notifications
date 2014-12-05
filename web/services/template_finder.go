package services

import (
	"fmt"
	"strings"

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

func (finder TemplateFinder) Find(templateName string) (models.Template, error) {
	names := finder.ParseTemplateName(templateName)
	if len(names) == 0 {
		return models.Template{}, TemplateNotFoundError(templateName)
	}
	template, err := finder.search(finder.database.Connection(), names[0], names[1:])
	if err != nil {
		return models.Template{}, err
	}

	return template, err
}

func (finder TemplateFinder) ParseTemplateName(name string) []string {
	names := make([]string, 0)

	for _, suffix := range models.TemplateNames {
		if strings.HasSuffix(name, suffix) {
			prefix := strings.TrimSuffix(name, suffix)
			prefix = strings.TrimSuffix(prefix, ".")
			var parts []string
			for _, part := range strings.Split(prefix, ".") {
				if part != "" {
					parts = append(parts, part)
				}
			}
			length := len(parts)
			if length > 2 {
				return names
			}
			for i := 0; i < length; i++ {
				beginning := strings.Join(parts, ".")
				names = append(names, beginning+"."+suffix)
				parts = parts[:len(parts)-1]
			}
			names = append(names, suffix)
			return names
		}
	}
	return names
}

func (finder TemplateFinder) search(connection models.ConnectionInterface, name string, alternates []string) (models.Template, error) {
	template, err := finder.templatesRepo.Find(connection, name)

	if err != nil {
		if (err == models.ErrRecordNotFound{}) {
			if len(alternates) > 0 {
				return finder.search(connection, alternates[0], alternates[1:])
			} else {
				return finder.findDefaultTemplate(name)
			}
		}
		return models.Template{}, err
	}

	return template, nil
}

func (finder TemplateFinder) findDefaultTemplate(name string) (models.Template, error) {
	if name == models.SubjectMissingTemplateName || name == models.SubjectProvidedTemplateName {
		return finder.defaultSubjectTemplate(name)
	} else {
		return finder.defaultTemplate(name)
	}
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

func (finder TemplateFinder) defaultSubjectTemplate(name string) (models.Template, error) {
	text, err := finder.readTemplate(name)
	if err != nil {
		return models.Template{}, err
	}
	return models.Template{Text: text, HTML: ""}, nil
}

func (finder TemplateFinder) readTemplate(name string) (string, error) {
	path := finder.rootPath + "/templates/" + name
	return finder.fileSystem.Read(path)
}
