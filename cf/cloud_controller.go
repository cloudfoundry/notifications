package cf

import (
	"fmt"

	"github.com/pivotal-cf-experimental/rainmaker"
)

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

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("CloudController Failure: %s", e.Message)
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
