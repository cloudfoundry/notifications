package handlers_test

import (
    "net/http"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
    fakes.RegisterFastTokenSigningMethod()

    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Suite")
}

type FakeNotify struct {
    Response []byte
    GUID     postal.TypedGUID
    Error    error
}

func (fake *FakeNotify) Execute(connection models.ConnectionInterface, req *http.Request,
    guid postal.TypedGUID, mailRecipe postal.MailRecipeInterface) ([]byte, error) {
    fake.GUID = guid

    return fake.Response, fake.Error
}
