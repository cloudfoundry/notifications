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
    TemplatesRepo models.TemplatesRepoInterface
    RootPath      string
    fileSystem    FileSystemInterface
    database      models.DatabaseInterface
}

type TemplateFinderInterface interface {
    Find(string) (models.Template, error)
}

func NewTemplateFinder(templatesRepo models.TemplatesRepoInterface, rootPath string, database models.DatabaseInterface, fileSystem FileSystemInterface) TemplateFinder {
    return TemplateFinder{
        TemplatesRepo: templatesRepo,
        RootPath:      rootPath,
        database:      database,
        fileSystem:    fileSystem,
    }
}

func (finder TemplateFinder) Find(templateName string) (models.Template, error) {
    var template models.Template
    var err error

    client, notificationType := parseTemplateName(templateName)
    templatesToSearchFor := []string{templateName, client + "." + notificationType, notificationType}
    connection := finder.database.Connection()

    for _, templateName := range templatesToSearchFor {
        template, err = finder.TemplatesRepo.Find(connection, templateName)
        if (template != models.Template{}) {
            break
        }
    }

    if (err == models.ErrRecordNotFound{}) {
        return finder.DefaultTemplate(notificationType)
    }

    return template, err
}

func (finder TemplateFinder) DefaultTemplate(notificationType string) (models.Template, error) {
    textPath := finder.RootPath + "/templates/" + notificationType + ".text"
    text, err := finder.fileSystem.Read(textPath)
    if err != nil {
        return models.Template{}, models.ErrRecordNotFound{}
    }

    htmlPath := finder.RootPath + "/templates/" + notificationType + ".html"
    html, err := finder.fileSystem.Read(htmlPath)
    if err != nil {
        return models.Template{}, models.ErrRecordNotFound{}
    }

    return models.Template{Text: text, HTML: html}, nil
}

func parseTemplateName(templateName string) (string, string) {
    client := ""
    notificationType := ""

    items := strings.Split(templateName, ".")
    numberOfItems := len(items)

    switch numberOfItems {
    case 3:
        client = items[0]
        notificationType = items[2]
    case 2:
        client = items[0]
        notificationType = items[1]
    case 1:
        notificationType = items[0]
    }
    return client, notificationType
}
