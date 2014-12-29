package support

import (
	"bytes"
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/web/params"
)

type TemplatesService struct {
	client *Client
}

type Template struct {
	Name     string                 `json:"name"`
	Subject  string                 `json:"subject"`
	Text     string                 `json:"text"`
	HTML     string                 `json:"html"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (t TemplatesService) Create(token string, template params.Template) (int, string, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, "", err
	}

	status, responseBody, err := t.client.makeRequest("POST", t.client.server.TemplatesBasePath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, "", err
	}

	var JSON struct {
		TemplateID string `json:"template_id"`
	}

	err = json.NewDecoder(responseBody).Decode(&JSON)
	if err != nil {
		return 0, "", err
	}

	return status, JSON.TemplateID, nil
}

func (t TemplatesService) Update(token string, id string, template params.Template) (int, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, err
	}

	status, _, err := t.client.makeRequest("PUT", t.client.server.TemplatePath(id), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (t TemplatesService) AssignToClient(token, clientID, templateID string) (int, error) {
	body, err := json.Marshal(map[string]string{
		"template": templateID,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := t.client.makeRequest("PUT", t.client.server.ClientsTemplatePath(clientID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (t TemplatesService) Get(token, templateID string) (int, Template, error) {
	var template Template

	status, body, err := t.client.makeRequest("GET", t.client.server.TemplatePath(templateID), nil, token)
	if err != nil {
		return 0, template, err
	}

	err = json.NewDecoder(body).Decode(&template)
	if err != nil {
		return 0, template, err
	}

	return status, template, nil
}
