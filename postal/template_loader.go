package postal

import "github.com/cloudfoundry-incubator/notifications/web/services"

const (
    SubjectMissingTemplateName  = "subject.missing"
    SubjectProvidedTemplateName = "subject.provided"
)

type Templates struct {
    Subject string
    Text    string
    HTML    string
}

type TemplatesLoaderInterface interface {
    LoadTemplates(string, string, string, string) (Templates, error)
}

type TemplatesLoader struct {
    finder services.TemplateFinderInterface
}

func NewTemplatesLoader(finder services.TemplateFinderInterface) TemplatesLoader {
    return TemplatesLoader{
        finder: finder,
    }
}

func (loader TemplatesLoader) LoadTemplates(subjectSuffix, contentSuffix, client, kind string) (Templates, error) {
    contentPath := client + "." + kind + "." + contentSuffix
    contentTemplate, err := loader.finder.Find(contentPath)
    if err != nil {
        return Templates{}, err
    }

    subjectPath := client + "." + kind + "." + subjectSuffix
    subjectTemplate, err := loader.finder.Find(subjectPath)
    if err != nil {
        return Templates{}, err
    }

    templates := Templates{
        Subject: subjectTemplate.Text,
        Text:    contentTemplate.Text,
        HTML:    contentTemplate.HTML,
    }

    return templates, nil
}
