package support

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CampaignType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
	TemplateID  string `json:"template_id"`
}

type CampaignTypesService struct {
	config Config
}

func NewCampaignTypesService(config Config) CampaignTypesService {
	return CampaignTypesService{
		config: config,
	}
}

func (c CampaignTypesService) Update(senderID, campaignTypeID, name, description, templateID string, critical bool, token string) (CampaignType, error) {
	var campaignType CampaignType

	content, err := json.Marshal(map[string]interface{}{
		"name":        name,
		"description": description,
		"critical":    critical,
		"template_id": templateID,
	})
	if err != nil {
		return campaignType, err
	}

	status, body, err := NewClient(c.config).makeRequest("PUT", c.config.Host+"/senders/"+senderID+"/campaign_types/"+campaignTypeID, bytes.NewBuffer(content), token)
	if err != nil {
		return campaignType, err
	}

	if status != http.StatusOK {
		return campaignType, UnexpectedStatusError{status, string(body)}
	}

	err = json.Unmarshal(body, &campaignType)
	if err != nil {
		return campaignType, err
	}

	return campaignType, nil
}

func (n CampaignTypesService) Create(senderID, name, description, templateID string, critical bool, token string) (CampaignType, error) {
	var campaignType CampaignType

	content, err := json.Marshal(map[string]interface{}{
		"name":        name,
		"description": description,
		"critical":    critical,
		"template_id": templateID,
	})
	if err != nil {
		return campaignType, err
	}

	status, body, err := NewClient(n.config).makeRequest("POST", n.config.Host+"/senders/"+senderID+"/campaign_types", bytes.NewBuffer(content), token)
	if err != nil {
		return campaignType, err
	}

	if status != http.StatusCreated {
		return campaignType, UnexpectedStatusError{status, string(body)}
	}

	err = json.Unmarshal(body, &campaignType)
	if err != nil {
		return campaignType, err
	}

	return campaignType, nil
}

func (n CampaignTypesService) Show(senderID, campaignTypeID, token string) (CampaignType, error) {
	var campaignType CampaignType

	status, body, err := NewClient(n.config).makeRequest("GET", n.config.Host+"/senders/"+senderID+"/campaign_types/"+campaignTypeID, nil, token)
	if err != nil {
		return campaignType, err
	}

	if status == http.StatusNotFound {
		return campaignType, NotFoundError{status, string(body)}
	} else if status != http.StatusOK {
		return campaignType, UnexpectedStatusError{status, string(body)}
	}

	err = json.Unmarshal(body, &campaignType)
	if err != nil {
		return campaignType, err
	}

	return campaignType, nil
}
