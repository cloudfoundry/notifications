package services

import (
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

func NewTemplateFinder(templatesRepo models.TemplatesRepoInterface, rootPath string, database models.DatabaseInterface, fileSystem FileSystemInterface) TemplateFinder {
    return TemplateFinder{
        templatesRepo: templatesRepo,
        rootPath:      rootPath,
        database:      database,
        fileSystem:    fileSystem,
    }
}

func (finder TemplateFinder) Find(templateName string) (models.Template, error) {
    templateNames := finder.templatesToSeachFor(templateName)
    template, err := finder.search(finder.database.Connection(), templateName, templateNames)
    if err != nil {
        return models.Template{}, err
    }

    if template.Name == templateName {
        template.Overridden = true
    }

    return template, err
}

func (finder TemplateFinder) templatesToSeachFor(templateName string) []string {
    var client string
    var notificationType string

    items := strings.Split(templateName, ".")
    numberOfItems := len(items)

    switch numberOfItems {
    case 4:
        client = items[0]
        notificationType = items[2] + "." + items[3]
    case 3:
        client = items[0]
        if items[1] == "subject" {
            notificationType = items[1] + "." + items[2]
        } else {
            notificationType = items[2]
        }
    case 2:
        client = items[0]
        notificationType = items[1]
    case 1:
        notificationType = items[0]
    }

    names := make([]string, 0)
    if client != "" {
        names = append(names, client+"."+notificationType)
    }

    names = append(names, notificationType)

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
