package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fileUtilities"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/dgrijalva/jwt-go"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type NotifyHelper struct {
    cloudController cf.CloudControllerInterface
    logger          *log.Logger
    uaaClient       uaa.UAAInterface
    guidGenerator   GUIDGenerationFunc
    mailClient      mail.ClientInterface
}

func NewNotifyHelper(cloudController cf.CloudControllerInterface, logger *log.Logger,
    uaaClient uaa.UAAInterface, guidGenerator GUIDGenerationFunc,
    mailClient mail.ClientInterface) NotifyHelper {
    return NotifyHelper{
        cloudController: cloudController,
        logger:          logger,
        uaaClient:       uaaClient,
        guidGenerator:   guidGenerator,
        mailClient:      mailClient,
    }
}

func (helper NotifyHelper) Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}

func (helper NotifyHelper) LoadUser(w http.ResponseWriter, guid string, uaaClient uaa.UAAInterface) (uaa.User, bool) {
    user, err := uaaClient.UserByID(guid)
    if err != nil {
        switch err.(type) {
        case *url.Error:
            w.WriteHeader(http.StatusBadGateway)
        case uaa.Failure:
            w.WriteHeader(http.StatusGone)
        default:
            w.WriteHeader(http.StatusInternalServerError)
        }
        return uaa.User{}, false
    }
    return user, true
}

func (helper NotifyHelper) SendMail(w http.ResponseWriter, req *http.Request,
    loadGuid func(path string) string,
    loadUsers func(spaceGuid, accessToken string) ([]cf.CloudControllerUser, error),
    loadSpace bool) {

    guid := loadGuid(req.URL.Path)

    params, err := NewNotifyParams(req.Body)

    if err != nil {
        helper.Error(w, 422, []string{"Request body could not be parsed"})
        return
    }

    if !params.Validate() {
        helper.Error(w, 422, params.Errors)
        return
    }

    token, err := helper.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    helper.uaaClient.SetToken(token.Access)

    ccUsers, err := loadUsers(guid, token.Access)
    if err != nil {
        helper.Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    env := config.NewEnvironment()

    var space, organization string
    if loadSpace {
        space, organization, err = helper.loadSpaceAndOrganization(guid, token.Access)
        if err != nil {
            helper.Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
            return
        }
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    responseInformation := make([]map[string]string, len(ccUsers))
    for index, ccUser := range ccUsers {
        helper.logger.Println(ccUser.Guid)
        user, ok := helper.LoadUser(w, ccUser.Guid, helper.uaaClient)
        if !ok {
            return
        }

        if len(user.Emails) > 0 {
            plainTextTemplate, htmlTemplate, err := helper.loadTemplates(loadSpace, env.RootPath)
            if err != nil {
                helper.Error(w, http.StatusInternalServerError, []string{"An email template could not be loaded"})
                return
            }
            context := helper.BuildSpaceContext(user, params, env, space, organization,
                clientToken.Claims["client_id"].(string), helper.guidGenerator,
                plainTextTemplate, htmlTemplate)
            status := helper.SendMailToUser(context, helper.logger, helper.mailClient)
            helper.logger.Println(status)

            userInfo := make(map[string]string)
            userInfo["status"] = status
            userInfo["recipient"] = ccUser.Guid
            userInfo["notification_id"] = context.MessageID

            responseInformation[index] = userInfo
        }
    }

    response := helper.generateResponse(responseInformation)
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}

func (helper NotifyHelper) loadTemplates(isSpace bool, rootPath string) (string, string, error) {
    var plainTextFileName, htmlFileName string
    if isSpace {
        plainTextFileName = rootPath + "/templates/space_body.text"
        htmlFileName = rootPath + "/templates/space_body.html"
    } else {
        plainTextFileName = rootPath + "/templates/user_body.text"
        htmlFileName = rootPath + "/templates/user_body.html"
    }

    plainTextTemplate, err := fileUtilities.ReadFile(plainTextFileName)
    if err != nil {
        return "", "", err
    }

    htmlTemplate, err := fileUtilities.ReadFile(htmlFileName)
    if err != nil {
        return "", "", err
    }

    return plainTextTemplate, htmlTemplate, nil
}

func (helper NotifyHelper) generateResponse(userInformation []map[string]string) []byte {
    response, err := json.Marshal(userInformation)
    if err != nil {
        panic(err)
    }

    return response
}

func (helper NotifyHelper) loadSpaceAndOrganization(spaceGuid, token string) (string, string, error) {
    space, err := helper.cloudController.LoadSpace(spaceGuid, token)
    if err != nil {
        return "", "", err
    }

    org, err := helper.cloudController.LoadOrganization(space.OrganizationGuid, token)
    if err != nil {
        return "", "", err
    }

    return space.Name, org.Name, nil
}

func (helper NotifyHelper) SendMailToUser(context MessageContext, logger *log.Logger, mailClient mail.ClientInterface) string {
    logger.Printf("Sending email to %s", context.To)
    status, message, err := SendMail(mailClient, context)
    if err != nil {
        panic(err)
    }

    logger.Print(message.Data())
    return status
}

func (helper NotifyHelper) BuildSpaceContext(user uaa.User, params NotifyParams, env config.Environment, space, organization, clientID string, guidGenerator GUIDGenerationFunc, plainTextEmailTemplate, htmlEmailTemplate string) MessageContext {
    return helper.buildContext(user, params, env, space, organization, clientID, guidGenerator, plainTextEmailTemplate, htmlEmailTemplate)
}

func (helper NotifyHelper) BuildUserContext(user uaa.User, params NotifyParams, env config.Environment, clientID string, guidGenerator GUIDGenerationFunc, plainTextEmailTemplate, htmlEmailTemplate string) MessageContext {
    return helper.buildContext(user, params, env, "", "", clientID, guidGenerator, plainTextEmailTemplate, htmlEmailTemplate)
}

func (handler NotifyHelper) buildContext(user uaa.User, params NotifyParams, env config.Environment, space, organization, clientID string, guidGenerator GUIDGenerationFunc, plainTextEmailTemplate, htmlEmailTemplate string) MessageContext {
    guid, err := guidGenerator()
    if err != nil {
        panic(err)
    }

    var kindDescription string
    if params.KindDescription == "" {
        kindDescription = params.Kind
    } else {
        kindDescription = params.KindDescription
    }

    var sourceDescription string
    if params.SourceDescription == "" {
        sourceDescription = clientID
    } else {
        sourceDescription = params.SourceDescription
    }

    return MessageContext{
        From:    env.Sender,
        To:      user.Emails[0],
        Subject: params.Subject,
        Text:    params.Text,
        HTML:    params.HTML,
        PlainTextEmailTemplate: plainTextEmailTemplate,
        HTMLEmailTemplate:      htmlEmailTemplate,
        KindDescription:        kindDescription,
        SourceDescription:      sourceDescription,
        ClientID:               clientID,
        MessageID:              guid.String(),
        Space:                  space,
        Organization:           organization,
    }
}
