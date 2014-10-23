package services

import (
    "fmt"
    "io/ioutil"

    "github.com/cloudfoundry-incubator/notifications/models"
)

const SpaceBody = "space_body"
const UserBody = "user_body"

type TemplateFinder struct {
    TemplatesRepo models.TemplatesRepoInterface
    RootPath      string
}

type TemplateFinderInterface interface {
    Find(string, string) (models.Template, error)
}

func NewTemplateFinder(templatesRepo models.TemplatesRepoInterface, rootPath string) TemplateFinder {
    return TemplateFinder{
        TemplatesRepo: templatesRepo,
        RootPath:      rootPath,
    }
}

func (finder TemplateFinder) Find(notificationType, templateName string) (models.Template, error) {
    template, err := finder.TemplatesRepo.Find(templateName)

    if (err == models.ErrRecordNotFound{}) {
        if notificationType == SpaceBody {
            return finder.DefaultTemplate(SpaceBody)
        } else {
            return finder.DefaultTemplate(UserBody)
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
