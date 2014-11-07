package fakes

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/models"
)

type FileSystem struct {
    Files map[string]string
}

func NewFileSystem(rootPath string) FileSystem {
    return FileSystem{
        Files: map[string]string{
            rootPath + "/templates/" + models.SpaceBodyTemplateName + ".text": "default-space-text",
            rootPath + "/templates/" + models.SpaceBodyTemplateName + ".html": "default-space-html",
            rootPath + "/templates/" + models.SubjectMissingTemplateName:      "default-missing-subject",
            rootPath + "/templates/" + models.SubjectProvidedTemplateName:     "default-provided-subject",
            rootPath + "/templates/" + models.UserBodyTemplateName + ".text":  "default-user-text",
            rootPath + "/templates/" + models.UserBodyTemplateName + ".html":  "default-user-html",
            rootPath + "/templates/" + models.EmailBodyTemplateName + ".html": "email-body-html",
            rootPath + "/templates/" + models.EmailBodyTemplateName + ".text": "email-body-text",
        },
    }
}

func (fs FileSystem) Exists(path string) bool {
    _, ok := fs.Files[path]
    return ok
}

func (fs FileSystem) Read(path string) (string, error) {
    if file, ok := fs.Files[path]; ok {
        return file, nil
    }
    return "", errors.New("File does not exist")
}
