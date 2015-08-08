package support

import (
	"bytes"
	"encoding/json"
)

type TemplatesService struct {
	client  *Client
	Default *DefaultTemplateService
}

func (s TemplatesService) Associations(token, templateID string) (int, []TemplateAssociation, error) {
	var associations TemplateAssociations

	status, body, err := s.client.makeRequest("GET", s.client.TemplateAssociationsPath(templateID), nil, token)
	if err != nil {
		return 0, associations.Associations, err
	}

	err = json.Unmarshal(body, &associations)
	if err != nil {
		return 0, associations.Associations, err
	}

	return status, associations.Associations, nil
}

func (s TemplatesService) Create(token string, template Template) (int, string, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, "", err
	}

	status, responseBody, err := s.client.makeRequest("POST", s.client.TemplatesPath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, "", err
	}

	var JSON struct {
		TemplateID string `json:"template_id"`
	}

	err = json.Unmarshal(responseBody, &JSON)
	if err != nil {
		return 0, "", err
	}

	return status, JSON.TemplateID, nil
}

func (s TemplatesService) Get(token, templateID string) (int, Template, error) {
	var template Template

	status, body, err := s.client.makeRequest("GET", s.client.TemplatePath(templateID), nil, token)
	if err != nil {
		return 0, template, err
	}

	err = json.Unmarshal(body, &template)
	if err != nil {
		return 0, template, err
	}

	return status, template, nil
}

func (s TemplatesService) Update(token, id string, template Template) (int, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PUT", s.client.TemplatePath(id), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (s TemplatesService) Delete(token, id string) (int, error) {
	status, _, err := s.client.makeRequest("DELETE", s.client.TemplatePath(id), nil, token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (s TemplatesService) List(token string) (int, []TemplateListItem, error) {
	var list []TemplateListItem
	status, body, err := s.client.makeRequest("GET", s.client.TemplatesPath(), nil, token)
	if err != nil {
		return 0, list, err
	}

	var templates map[string]TemplateListItem
	err = json.Unmarshal(body, &templates)
	if err != nil {
		panic(err)
	}

	for id, template := range templates {
		template.ID = id
		list = append(list, template)
	}

	return status, list, err
}

func (s TemplatesService) AssignToClient(token, clientID, templateID string) (int, error) {
	body, err := json.Marshal(map[string]string{
		"template": templateID,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PUT", s.client.ClientsTemplatePath(clientID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (s TemplatesService) AssignToNotification(token, clientID, notificationID, templateID string) (int, error) {
	body, err := json.Marshal(map[string]string{
		"template": templateID,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PUT", s.client.ClientsNotificationsTemplatePath(clientID, notificationID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
