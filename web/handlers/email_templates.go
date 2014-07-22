package handlers

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fileUtilities"
)

type EmailTemplateManager struct {
    ReadFile   func(string) (string, error)
    FileExists func(string) bool
}

func NewTemplateManager() EmailTemplateManager {
    return EmailTemplateManager{
        ReadFile:   fileUtilities.ReadFile,
        FileExists: fileUtilities.FileExists,
    }
}

func (manager EmailTemplateManager) LoadEmailTemplate(filename string) (string, error) {
    env := config.NewEnvironment()
    templatesDirectory := "/templates"

    basePath := env.RootPath + templatesDirectory
    defaultPath := basePath + "/" + filename
    overRidePath := basePath + "/overrides/" + filename

    if manager.FileExists(overRidePath) {
        fileContents, err := manager.ReadFile(overRidePath)
        if err != nil {
            return "", err
        }
        return fileContents, nil
    }

    fileContents, err := manager.ReadFile(defaultPath)
    if err != nil {
        return "", err
    }
    return fileContents, nil
}
