package support

import (
	"bytes"
	"encoding/json"
)

type DefaultTemplateService struct {
	client *Client
}

func (s DefaultTemplateService) Get(token string) (int, Template, error) {
	var template Template

	status, body, err := s.client.makeRequest("GET", s.client.DefaultTemplatePath(), nil, token)
	if err != nil {
		return 0, template, err
	}

	err = json.NewDecoder(body).Decode(&template)
	if err != nil {
		return 0, template, err
	}

	return status, template, nil
}

func (s DefaultTemplateService) Update(token string, template Template) (int, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PUT", s.client.DefaultTemplatePath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
