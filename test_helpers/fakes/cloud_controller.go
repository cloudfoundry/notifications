package fakes

import (
    "fmt"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/cf"
)

type FakeCloudController struct {
    CurrentToken             string
    GetUsersBySpaceGuidError error
    LoadSpaceError           error
    LoadOrganizationError    error
    UsersBySpaceGuid         map[string][]cf.CloudControllerUser
    Spaces                   map[string]cf.CloudControllerSpace
    Orgs                     map[string]cf.CloudControllerOrganization
}

func NewFakeCloudController() *FakeCloudController {
    return &FakeCloudController{
        UsersBySpaceGuid: make(map[string][]cf.CloudControllerUser),
    }
}

func (fake *FakeCloudController) GetUsersBySpaceGuid(guid, token string) ([]cf.CloudControllerUser, error) {
    fake.CurrentToken = token

    if users, ok := fake.UsersBySpaceGuid[guid]; ok {
        return users, fake.GetUsersBySpaceGuidError
    } else {
        return make([]cf.CloudControllerUser, 0), fake.GetUsersBySpaceGuidError
    }
}

func (fake *FakeCloudController) LoadSpace(guid, token string) (cf.CloudControllerSpace, error) {
    if fake.LoadSpaceError != nil {
        return cf.CloudControllerSpace{}, fake.LoadSpaceError
    }

    if space, ok := fake.Spaces[guid]; ok {
        return space, nil
    } else {
        return cf.CloudControllerSpace{}, cf.NewFailure(http.StatusNotFound, fmt.Sprintf(`{"code":40004,"description":"The app space could not be found: %s","error_code":"CF-SpaceNotFound"}`, guid))
    }
}

func (fake *FakeCloudController) LoadOrganization(guid, token string) (cf.CloudControllerOrganization, error) {
    if fake.LoadOrganizationError != nil {
        return cf.CloudControllerOrganization{}, fake.LoadOrganizationError
    }

    if org, ok := fake.Orgs[guid]; ok {
        return org, nil
    } else {
        return cf.CloudControllerOrganization{}, cf.NewFailure(http.StatusNotFound, fmt.Sprintf(`{"code":30003,"description":"The organization could not be found: %s","error_code":"CF-OrganizationNotFound"}`, guid))
    }
}
