package mocks

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/docs"
)

func NewRequestInspector() *RequestInspector {
	return &RequestInspector{}
}

type RequestInspector struct {
	GetResourceInfoCall struct {
		Receives struct {
			Request *http.Request
		}

		Returns struct {
			ResourceInfo docs.ResourceInfo
		}
	}
}

func (r *RequestInspector) GetResourceInfo(request *http.Request) docs.ResourceInfo {
	r.GetResourceInfoCall.Receives.Request = request
	return r.GetResourceInfoCall.Returns.ResourceInfo
}
