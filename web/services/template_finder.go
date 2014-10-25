package services

import (
    "fmt"
    "io/ioutil"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
)

const (
    SpaceBody = "space_body"
    UserBody  = "user_body"
    EmailBody = "email_body"
)

type TemplateFinder struct {
    TemplatesRepo models.TemplatesRepoInterface
    RootPath      string
    database      models.DatabaseInterface
}

type TemplateFinderInterface interface {
    Find(string) (models.Template, error)
}

func NewTemplateFinder(templatesRepo models.TemplatesRepoInterface, rootPath string, database models.DatabaseInterface) TemplateFinder {
    return TemplateFinder{
        TemplatesRepo: templatesRepo,
        RootPath:      rootPath,
        database:      database,
    }
}

func (finder TemplateFinder) Find(templateName string) (models.Template, error) {
    var template models.Template
    var err error
    notificationType := parseNotificationType(templateName)
    client := parseClientType(templateName)
    templatesToSearchFor := []string{templateName, client + "." + notificationType, notificationType}
    connection := finder.database.Connection()

    for _, templateName := range templatesToSearchFor {
        template, err = finder.TemplatesRepo.Find(connection, templateName)
        if (template != models.Template{}) {
            break
        }
    }

    if (err == models.ErrRecordNotFound{}) {
        switch notificationType {
        case SpaceBody:
            return finder.DefaultTemplate(SpaceBody)
        case UserBody:
            return finder.DefaultTemplate(UserBody)
        case EmailBody:
            return finder.DefaultTemplate(EmailBody)
        }
    }

    return template, err
}

func (finder TemplateFinder) DefaultTemplate(notificationType string) (models.Template, error) {
    textPath := finder.RootPath + "/templates/" + notificationType + ".text"
    bytes, err := ioutil.ReadFile(textPath)
    if err != nil {
        return models.Template{}, fmt.Errorf("Could not read text file")
    }
    text := string(bytes)

    htmlPath := finder.RootPath + "/templates/" + notificationType + ".html"
    bytes, err = ioutil.ReadFile(htmlPath)
    if err != nil {
        return models.Template{}, fmt.Errorf("Could not read html file")
    }
    html := string(bytes)

    return models.Template{Text: text, HTML: html}, nil
}

func parseClientType(templateName string) string {
    theSplit := strings.Split(templateName, ".")
    if len(theSplit) == 3 || len(theSplit) == 2 {
        return theSplit[0]
    } else {
        return ""
    }
}

func parseNotificationType(templateName string) string {
    if strings.HasSuffix(templateName, UserBody) {
        return UserBody
    } else if strings.HasSuffix(templateName, SpaceBody) {
        return SpaceBody
    } else if strings.HasSuffix(templateName, EmailBody) {
        return EmailBody
    }
    return ""
}
