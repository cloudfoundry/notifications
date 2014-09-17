package handlers_test

import (
    "bytes"
    "log"
    "net/http"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
    fakes.RegisterFastTokenSigningMethod()

    buffer := bytes.NewBuffer([]byte{})
    metricsLogger := metrics.Logger
    metrics.Logger = log.New(buffer, "", 0)

    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Suite")

    metrics.Logger = metricsLogger
}

type FakeNotify struct {
    Response []byte
    GUID     postal.TypedGUID
    Error    error
}

func (fake *FakeNotify) Execute(connection models.ConnectionInterface, req *http.Request, context stack.Context,
    guid postal.TypedGUID, mailRecipe postal.MailRecipeInterface) ([]byte, error) {
    fake.GUID = guid

    return fake.Response, fake.Error
}
