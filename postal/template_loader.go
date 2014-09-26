package postal

import "github.com/cloudfoundry-incubator/notifications/config"

const (
    SubjectMissingTemplateName  = "subject.missing"
    SubjectProvidedTemplateName = "subject.provided"
)

type Templates struct {
    Subject  string
    Text     string
    HTML     string
    UserGUID string
}

type TemplateLoaderInterface interface {
    LoadNamedTemplates(string, string, string) (Templates, error)
    LoadNamedTemplatesWithClientAndKind(string, string, string, string, string) (Templates, error)
}

type TemplateLoader struct {
    fs       FileSystemInterface
    ClientID string
    Kind     string
}

func NewTemplateLoader(fs FileSystemInterface) TemplateLoader {
    return TemplateLoader{
        fs: fs,
    }
}

func (loader TemplateLoader) LoadNamedTemplates(subjectTemplateName, textTemplateName, htmlTemplateName string) (Templates, error) {
    var err error
    templates := Templates{}

    templates.Subject, err = loader.LoadTemplate(subjectTemplateName)
    if err != nil {
        return templates, err
    }

    templates.Text, err = loader.LoadTemplate(textTemplateName)
    if err != nil {
        return templates, err
    }

    templates.HTML, err = loader.LoadTemplate(htmlTemplateName)
    if err != nil {
        return templates, err
    }
    return templates, nil
}

func (loader TemplateLoader) LoadNamedTemplatesWithClientAndKind(subjectTemplateName, textTemplateName, htmlTemplateName, clientID, kind string) (Templates, error) {
    var err error
    templates := Templates{}
    loader.ClientID = clientID
    loader.Kind = kind

    templates.Subject, err = loader.LoadTemplate(subjectTemplateName)
    if err != nil {
        return templates, err
    }

    templates.Text, err = loader.LoadTemplate(textTemplateName)
    if err != nil {
        return templates, err
    }

    templates.HTML, err = loader.LoadTemplate(htmlTemplateName)
    if err != nil {
        return templates, err
    }
    return templates, nil
}

func (loader TemplateLoader) loadSubject(subject string) (string, error) {
    if subject == "" {
        return loader.LoadTemplate(SubjectMissingTemplateName)
    } else {
        return loader.LoadTemplate(SubjectProvidedTemplateName)
    }
}

func (loader TemplateLoader) LoadTemplate(filename string) (string, error) {
    env := config.NewEnvironment()

    clientKindOverridePath := env.RootPath + "/templates/overrides/" + loader.ClientID + "." + loader.Kind + "." + filename
    if loader.fs.Exists(clientKindOverridePath) {
        return loader.fs.Read(clientKindOverridePath)
    }

    clientOverridePath := env.RootPath + "/templates/overrides/" + loader.ClientID + "." + filename
    if loader.fs.Exists(clientOverridePath) {
        return loader.fs.Read(clientOverridePath)
    }

    overRidePath := env.RootPath + "/templates/overrides/" + filename
    if loader.fs.Exists(overRidePath) {
        return loader.fs.Read(overRidePath)
    }

    return loader.fs.Read(env.RootPath + "/templates/" + filename)
}
