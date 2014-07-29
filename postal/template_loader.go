package postal

import "github.com/cloudfoundry-incubator/notifications/config"

const (
    SubjectMissingTemplateName  = "subject.missing"
    SubjectProvidedTemplateName = "subject.provided"
    SpaceTextTemplateName       = "space_body.text"
    SpaceHTMLTemplateName       = "space_body.html"
    UserTextTemplateName        = "user_body.text"
    UserHTMLTemplateName        = "user_body.html"
)

type Templates struct {
    Subject string
    Text    string
    HTML    string
}

type TemplateLoader struct {
    fs FileSystemInterface
}

func NewTemplateLoader(fs FileSystemInterface) TemplateLoader {
    return TemplateLoader{
        fs: fs,
    }
}

func (loader TemplateLoader) Load(subject string, notificationType NotificationType) (Templates, error) {
    var err error
    templates := Templates{}

    templates.Subject, err = loader.loadSubject(subject)
    if err != nil {
        return templates, err
    }

    templates.Text, err = loader.loadText(notificationType)
    if err != nil {
        return templates, err
    }

    templates.HTML, err = loader.loadHTML(notificationType)
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

func (loader TemplateLoader) loadText(notificationType NotificationType) (string, error) {
    if notificationType == IsSpace {
        return loader.LoadTemplate(SpaceTextTemplateName)
    } else {
        return loader.LoadTemplate(UserTextTemplateName)
    }
}

func (loader TemplateLoader) loadHTML(notificationType NotificationType) (string, error) {
    if notificationType == IsSpace {
        return loader.LoadTemplate(SpaceHTMLTemplateName)
    } else {
        return loader.LoadTemplate(UserHTMLTemplateName)
    }

}

func (loader TemplateLoader) LoadTemplate(filename string) (string, error) {
    env := config.NewEnvironment()

    overRidePath := env.RootPath + "/templates/overrides/" + filename
    if loader.fs.Exists(overRidePath) {
        return loader.fs.Read(overRidePath)
    }

    return loader.fs.Read(env.RootPath + "/templates/" + filename)
}
