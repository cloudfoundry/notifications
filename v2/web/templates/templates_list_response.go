package templates

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type TemplatesListResponseLinks struct {
	Self Link `json:"self"`
}

type TemplatesListResponse struct {
	Templates []TemplateResponse         `json:"templates"`
	Links     TemplatesListResponseLinks `json:"_links"`
}

func NewTemplatesListResponse(templateList []collections.Template) TemplatesListResponse {
	templates := []TemplateResponse{}

	for _, t := range templateList {
		templates = append(templates, NewTemplateResponse(t))
	}

	return TemplatesListResponse{
		Templates: templates,
		Links:     TemplatesListResponseLinks{Link{"/templates"}},
	}
}
