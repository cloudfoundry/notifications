package templates

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type Link struct {
	Href string `json:"href"`
}

type TemplateResponseLinks struct {
	Self Link `json:"self"`
}

type TemplateResponse struct {
	ID       string                `json:"id"`
	Name     string                `json:"name"`
	Text     string                `json:"text"`
	HTML     string                `json:"html"`
	Subject  string                `json:"subject"`
	Metadata *json.RawMessage      `json:"metadata"`
	Links    TemplateResponseLinks `json:"_links"`
}

func NewTemplateResponse(template collections.Template) TemplateResponse {
	metadata := json.RawMessage(template.Metadata)
	return TemplateResponse{
		ID:       template.ID,
		Name:     template.Name,
		Text:     template.Text,
		HTML:     template.HTML,
		Subject:  template.Subject,
		Metadata: &metadata,
		Links:    TemplateResponseLinks{Link{fmt.Sprintf("/templates/%s", template.ID)}},
	}
}
