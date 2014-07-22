package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
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

func Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(NotifyFailureResponse{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}

func (helper NotifyHelper) NotifyServeHTTP(w http.ResponseWriter, req *http.Request,
    guid string, loadCCUsers func(spaceGuid, accessToken string) ([]cf.CloudControllerUser, error),
    loadSpace bool) {

    params, err := NewNotifyParams(req.Body)
    if err != nil {
        Error(w, 422, []string{"Request body could not be parsed"})
        return
    }

    if !params.Validate() {
        Error(w, 422, params.Errors)
        return
    }

    token, err := helper.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    helper.uaaClient.SetToken(token.Access)

    ccUsers, err := loadCCUsers(guid, token.Access)
    if err != nil {
        Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    uaaUsers := make([]uaa.User, len(ccUsers))
    for index, ccUser := range ccUsers {
        helper.logger.Println(ccUser.Guid)
        user, ok := helper.LoadUaaUser(w, ccUser.Guid, helper.uaaClient)
        if !ok {
            return
        }
        uaaUsers[index] = user
    }

    var space, organization string
    if loadSpace {
        space, organization, err = helper.loadSpaceAndOrganization(guid, token.Access)
        if err != nil {
            Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
            return
        }
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    responseGenerator := NewNotifyResponseGenerator(helper.logger, helper.guidGenerator,
        helper.mailClient, helper.uaaClient)

    responseGenerator.GenerateResponse(uaaUsers, params, space,
        organization, clientToken.Claims["client_id"].(string), w, loadSpace)
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

func (helper NotifyHelper) LoadUaaUser(w http.ResponseWriter, guid string, uaaClient uaa.UAAInterface) (uaa.User, bool) {
    user, err := uaaClient.UserByID(guid)
    if err != nil {
        switch err.(type) {
        case *url.Error:
            Error(w, http.StatusBadGateway, []string{"UAA is unavailable"})
        case uaa.Failure:
            Error(w, http.StatusGone, []string{"UAA is unavailable"})
        default:
            Error(w, http.StatusInternalServerError, []string{"UAA is unavailable"})
        }
        return uaa.User{}, false
    }
    return user, true
}
