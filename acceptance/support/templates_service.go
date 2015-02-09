package support

import (
	"bytes"
	"encoding/json"
)

type TemplatesService struct {
	client  *Client
	Default *DefaultTemplateService
}

func (t TemplatesService) Associations(token, templateID string) (int, []TemplateAssociation, error) {
	var associations TemplateAssociations

	status, body, err := t.client.makeRequest("GET", t.client.TemplateAssociationsPath(templateID), nil, token)
	if err != nil {
		return 0, associations.Associations, err
	}

	err = json.NewDecoder(body).Decode(&associations)
	if err != nil {
		return 0, associations.Associations, err
	}

	return status, associations.Associations, nil
}

func (t TemplatesService) Create(token string, template Template) (int, string, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, "", err
	}

	status, responseBody, err := t.client.makeRequest("POST", t.client.TemplatesPath(), bytes.NewBuffer(body), token)
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

func (t TemplatesService) Get(token, templateID string) (int, Template, error) {
	var template Template

	status, body, err := t.client.makeRequest("GET", t.client.TemplatePath(templateID), nil, token)
	if err != nil {
		return 0, template, err
	}

	err = json.NewDecoder(body).Decode(&template)
	if err != nil {
		return 0, template, err
	}

	return status, template, nil
}

func (t TemplatesService) Update(token, id string, template Template) (int, error) {
	body, err := json.Marshal(template)
	if err != nil {
		return 0, err
	}

	status, _, err := t.client.makeRequest("PUT", t.client.TemplatePath(id), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (t TemplatesService) Delete(token, id string) (int, error) {
	status, _, err := t.client.makeRequest("DELETE", t.client.TemplatePath(id), nil, token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (t TemplatesService) List(token string) (int, []TemplateListItem, error) {
	var list []TemplateListItem
	status, body, err := t.client.makeRequest("GET", t.client.TemplatesPath(), nil, token)
	if err != nil {
		return 0, list, err
	}

	var templates map[string]TemplateListItem
	err = json.NewDecoder(body).Decode(&templates)
	if err != nil {
		panic(err)
	}

	for id, template := range templates {
		template.ID = id
		list = append(list, template)
	}

	return status, list, err
}

func (t TemplatesService) AssignToClient(token, clientID, templateID string) (int, error) {
	body, err := json.Marshal(map[string]string{
		"template": templateID,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := t.client.makeRequest("PUT", t.client.ClientsTemplatePath(clientID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (t TemplatesService) AssignToNotification(token, clientID, notificationID, templateID string) (int, error) {
	body, err := json.Marshal(map[string]string{
		"template": templateID,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := t.client.makeRequest("PUT", t.client.ClientsNotificationsTemplatePath(clientID, notificationID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
