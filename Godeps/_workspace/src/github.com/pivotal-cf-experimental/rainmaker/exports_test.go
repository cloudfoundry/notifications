package rainmaker

import (
	"net/url"

	"github.com/pivotal-cf-experimental/rainmaker/internal/documents"
)

func NewRequestArguments(method, path, token string, body interface{}, statusCodes []int) requestArguments {
	return requestArguments{
		Method: method,
		Path:   path,
		Token:  token,
		Body:   body,
		AcceptableStatusCodes: statusCodes,
	}
}

func (client Client) MakeRequest(requestArgs requestArguments) (int, []byte, error) {
	return client.makeRequest(requestArgs)
}

func (client Client) Unmarshal(body []byte, response interface{}) error {
	return client.unmarshal(body, response)
}

func NewRequestPlan(path string, query url.Values) requestPlan {
	return newRequestPlan(path, query)
}

func NewOrganizationFromResponse(config Config, document documents.OrganizationResponse) Organization {
	return newOrganizationFromResponse(config, document)
}

func NewSpaceFromResponse(config Config, document documents.SpaceResponse) Space {
	return newSpaceFromResponse(config, document)
}
