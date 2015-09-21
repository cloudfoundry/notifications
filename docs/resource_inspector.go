package docs

import (
	"net/http"
	"regexp"
	"strings"
)

type RequestInspector struct {
}

type ResourceInfo struct {
	ResourceType string
	ListName     string
	ItemName     string
	IsItem       bool
}

func NewRequestInspector() *RequestInspector {
	return &RequestInspector{}
}

func capitalizeFirstLetter(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func (i *RequestInspector) GetResourceInfo(request *http.Request) ResourceInfo {
	var resourceInfo ResourceInfo

	parts := strings.Split(request.URL.Path, "/")

	resourceInfo.ResourceType = parts[len(parts)-1]

	guidValidator := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}")
	if guidValidator.Match([]byte(resourceInfo.ResourceType)) {
		resourceInfo.ResourceType = parts[len(parts)-2]
		resourceInfo.IsItem = true
	}

	resourceInfo.ListName = capitalizeFirstLetter(strings.Replace(resourceInfo.ResourceType, "_", " ", -1))
	resourceInfo.ItemName = resourceInfo.ListName[0 : len(resourceInfo.ListName)-1]

	return resourceInfo
}
