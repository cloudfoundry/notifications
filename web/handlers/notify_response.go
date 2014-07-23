package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type NotifyResponseGenerator struct {
    logger        *log.Logger
    guidGenerator GUIDGenerationFunc
    mailClient    mail.ClientInterface
}

func NewNotifyResponseGenerator(logger *log.Logger,
    guidGenerator GUIDGenerationFunc, mailClient mail.ClientInterface) NotifyResponseGenerator {
    return NotifyResponseGenerator{
        logger:        logger,
        guidGenerator: guidGenerator,
        mailClient:    mailClient,
    }
}

type NotifyResponse []map[string]string
type NotifyFailureResponse map[string][]string

func (generator NotifyResponseGenerator) GenerateResponse(uaaUsers []uaa.User,
    params NotifyParams, space, organization, token string,
    w http.ResponseWriter, loadSpace bool) {

    env := config.NewEnvironment()
    messages := make(NotifyResponse, len(uaaUsers))

    subjectTemplate, err := generator.LoadSubjectTemplate(params.Subject, NewTemplateManager())
    if err != nil {
        Error(w, http.StatusInternalServerError, []string{"An email template could not be loaded"})
        return
    }

    plainTextTemplate, htmlTemplate, err := generator.LoadBodyTemplates(loadSpace, NewTemplateManager())
    if err != nil {
        Error(w, http.StatusInternalServerError, []string{"An email template could not be loaded"})
        return
    }

    for index, uaaUser := range uaaUsers {
        if len(uaaUser.Emails) > 0 {
            context := NewMessageContext(uaaUser, params, env, space, organization,
                token, generator.guidGenerator, plainTextTemplate, htmlTemplate, subjectTemplate)

            emailStatus := generator.SendMailToUser(context, generator.logger, generator.mailClient)
            generator.logger.Println(emailStatus)

            mailInfo := make(map[string]string)
            mailInfo["status"] = emailStatus
            mailInfo["recipient"] = uaaUser.ID
            mailInfo["notification_id"] = context.MessageID

            messages[index] = mailInfo
        }
    }

    responseBytes, err := json.Marshal(messages)
    if err != nil {
        panic(err)
    }
    w.WriteHeader(http.StatusOK)
    w.Write(responseBytes)
}

func (generator NotifyResponseGenerator) SendMailToUser(context MessageContext, logger *log.Logger,
    mailClient mail.ClientInterface) string {

    logger.Printf("Sending email to %s", context.To)
    status, message, err := SendMail(mailClient, context)
    if err != nil {
        panic(err)
    }

    logger.Print(message.Data())
    return status
}

func (generator NotifyResponseGenerator) LoadSubjectTemplate(subject string, templateManager EmailTemplateManager) (string, error) {
    var templateToLoad string
    if subject == "" {
        templateToLoad = "subject.missing"
    } else {
        templateToLoad = "subject.provided"
    }

    subjectTemplate, err := templateManager.LoadEmailTemplate(templateToLoad)
    if err != nil {
        return "", err
    }

    return subjectTemplate, nil
}

func (generator NotifyResponseGenerator) LoadBodyTemplates(isSpace bool, templateManager EmailTemplateManager) (string, string, error) {
    var plainTextTemplate, htmlTemplate string
    var plainErr, htmlErr error

    if isSpace {
        plainTextTemplate, plainErr = templateManager.LoadEmailTemplate("space_body.text")
        htmlTemplate, htmlErr = templateManager.LoadEmailTemplate("space_body.html")
    } else {
        plainTextTemplate, plainErr = templateManager.LoadEmailTemplate("user_body.text")
        htmlTemplate, htmlErr = templateManager.LoadEmailTemplate("user_body.html")
    }

    if plainErr != nil {
        return "", "", plainErr
    }

    if htmlErr != nil {
        return "", "", htmlErr
    }

    return plainTextTemplate, htmlTemplate, nil
}
