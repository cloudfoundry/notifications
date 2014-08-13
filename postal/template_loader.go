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
    fs       FileSystemInterface
    ClientID string
    Kind     string
}

func NewTemplateLoader(fs FileSystemInterface) TemplateLoader {
    return TemplateLoader{
        fs: fs,
    }
}

func (loader TemplateLoader) Load(subject string, guid TypedGUID, clientID string, kind string) (Templates, error) {
    var err error
    loader.ClientID = clientID
    loader.Kind = kind
    templates := Templates{}

    templates.Subject, err = loader.loadSubject(subject)
    if err != nil {
        return templates, err
    }

    templates.Text, err = loader.loadText(guid)
    if err != nil {
        return templates, err
    }

    templates.HTML, err = loader.loadHTML(guid)
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

func (loader TemplateLoader) loadText(guid TypedGUID) (string, error) {
    if guid.BelongsToSpace() {
        return loader.LoadTemplate(SpaceTextTemplateName)
    } else {
        return loader.LoadTemplate(UserTextTemplateName)
    }
}

func (loader TemplateLoader) loadHTML(guid TypedGUID) (string, error) {
    if guid.BelongsToSpace() {
        return loader.LoadTemplate(SpaceHTMLTemplateName)
    } else {
        return loader.LoadTemplate(UserHTMLTemplateName)
    }

}

func (loader TemplateLoader) LoadTemplate(filename string) (string, error) {
    env := config.NewEnvironment()

    clientKindOverridePath := env.RootPath + "/templates/overrides/" + loader.ClientID + "." + loader.Kind + "." + filename
    if loader.fs.Exists(clientKindOverridePath) {
        return loader.fs.Read(clientKindOverridePath)
    }

    overRidePath := env.RootPath + "/templates/overrides/" + filename
    if loader.fs.Exists(overRidePath) {
        return loader.fs.Read(overRidePath)
    }

    return loader.fs.Read(env.RootPath + "/templates/" + filename)
}
