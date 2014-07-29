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

type UAAInterface interface {
    uaa.GetClientTokenInterface
    uaa.SetTokenInterface
    uaa.UsersByIDsInterface
}

type UAADownError struct {
    message string
}

func (err UAADownError) Error() string {
    return err.message
}

type UAAUserNotFoundError struct {
    message string
}

func (err UAAUserNotFoundError) Error() string {
    return err.message
}

type UAAGenericError struct {
    message string
}

func (err UAAGenericError) Error() string {
    return err.message
}

type NotifyHelper struct {
    cloudController cf.CloudControllerInterface
    logger          *log.Logger
    uaaClient       UAAInterface
    guidGenerator   GUIDGenerationFunc
    mailClient      mail.ClientInterface
}

func NewNotifyHelper(cloudController cf.CloudControllerInterface, logger *log.Logger,
    uaaClient UAAInterface, guidGenerator GUIDGenerationFunc,
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

    uaaUsers := make(map[string]uaa.User)
    var ccGuids []string
    for _, ccUser := range ccUsers {
        helper.logger.Println("CloudController user guid: " + ccUser.Guid)
        ccGuids = append(ccGuids, ccUser.Guid)
    }

    users, err := helper.LoadUaaUsers(ccGuids, helper.uaaClient)
    if err != nil {
        switch err.(type) {
        case UAADownError:
            Error(w, http.StatusBadGateway, []string{"UAA is unavailable"})
            return
        case UAAGenericError:
            Error(w, http.StatusBadGateway, []string{err.Error()})
            return
        default:
            Error(w, http.StatusBadGateway, []string{err.Error()})
            return
        }
    }

    for _, user := range users {
        uaaUsers[user.ID] = user
    }

    for _, ccUser := range ccUsers {
        if _, ok := uaaUsers[ccUser.Guid]; !ok {
            uaaUsers[ccUser.Guid] = uaa.User{}
        }
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
        helper.mailClient)

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

func (helper NotifyHelper) LoadUaaUsers(guids []string, uaaClient UAAInterface) ([]uaa.User, error) {
    users, err := uaaClient.UsersByIDs(guids...)
    if err != nil {
        switch err.(type) {
        case *url.Error:
            return users, UAADownError{
                message: "UAA is unavailable",
            }
        case uaa.Failure:
            uaaFailure := err.(uaa.Failure)
            helper.logger.Printf("error:  %v", err)

            if uaaFailure.Code() == 404 {
                if strings.Contains(uaaFailure.Message(), "Requested route") {
                    return users, UAADownError{
                        message: "UAA is unavailable",
                    }
                } else {
                    return users, UAAGenericError{
                        message: "UAA Unknown 404 error message: " + uaaFailure.Message(),
                    }
                }
            }

            return users, UAADownError{
                message: "UAA is unavailable",
            }
        default:
            return users, UAAGenericError{
                message: "UAA Unknown Error: " + err.Error(),
            }
        }
    }
    return users, nil
}
