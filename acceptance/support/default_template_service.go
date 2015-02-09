package support

import (
	"bytes"
	"encoding/json"
)

type DefaultTemplateService struct {
	client *Client
}

func (d DefaultTemplateService) Get(token string) (int, Template, error) {
	var template Template

	status, body, err := d.client.makeRequest("GET", d.client.DefaultTemplatePath(), nil, token)
	if err != nil {
		return 0, template, err
	}

	err = json.NewDecoder(body).Decode(&template)
	if err != nil {
		return 0, template, err
	}

	return status, template, nil
}

func (d DefaultTemplateService) Update(token string, template Template) (int, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, err
	}

	status, _, err := d.client.makeRequest("PUT", d.client.DefaultTemplatePath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
