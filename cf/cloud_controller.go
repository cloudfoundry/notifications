package cf

import (
	"fmt"

	"github.com/pivotal-golang/rainmaker"
)

type CloudControllerInterface interface {
	GetUsersBySpaceGuid(string, string) ([]CloudControllerUser, error)
	GetUsersByOrgGuid(string, string) ([]CloudControllerUser, error)
	GetManagersByOrgGuid(string, string) ([]CloudControllerUser, error)
	GetAuditorsByOrgGuid(string, string) ([]CloudControllerUser, error)
	GetBillingManagersByOrgGuid(string, string) ([]CloudControllerUser, error)
	LoadSpace(string, string) (CloudControllerSpace, error)
	LoadOrganization(string, string) (CloudControllerOrganization, error)
}

type CloudController struct {
	client rainmaker.Client
}

func NewCloudController(host string, skipVerifySSL bool) CloudController {
	return CloudController{
		client: rainmaker.NewClient(rainmaker.Config{
			Host:          host,
			SkipVerifySSL: skipVerifySSL,
		}),
	}
}

type CloudControllerUser struct {
	GUID string
}

type CloudControllerSpace struct {
	GUID             string
	Name             string
	OrganizationGUID string
}

type CloudControllerOrganization struct {
	GUID string
	Name string
}

type Failure struct {
	Code    int
	Message string
}

func NewFailure(code int, message string) Failure {
	return Failure{
		Code:    code,
		Message: message,
	}
}

func (failure Failure) Error() string {
	return fmt.Sprintf("CloudController Failure (%d): %s", failure.Code, failure.Message)
}
