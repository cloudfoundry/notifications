package params

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type Template struct {
	Name    string `json:"name"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
	Subject string `json:"subject"`
}

type TemplateCreateError struct{}

func (err TemplateCreateError) Error() string {
	return "Failed to create Template in the database"
}

func NewTemplate(body io.Reader) (Template, error) {
	var template Template

	jsonBody, err := ioutil.ReadAll(body)
	if err != nil {
		return Template{}, ParseError{}
	}

	err = json.Unmarshal(jsonBody, &template)
	if err != nil {
		return template, ParseError{}
	}

	err = containsArguments(template)
	if err != nil {
		return Template{}, err
	}

	template.setDefaults()

	return template, nil
}

func containsArguments(template Template) error {
	if template.Name == "" {
		return ValidationError([]string{"Request is missing the required field: name"})
	}
	if template.HTML == "" {
		return ValidationError([]string{"Request is missing the required field: html"})
	}

	return nil
}

func (t *Template) ToModel() models.Template {
	return models.Template{
		Name:    t.Name,
		Text:    t.Text,
		HTML:    t.HTML,
		Subject: t.Subject,
	}
}

func (template *Template) setDefaults() {
	if template.Subject == "" {
		template.Subject = "{{.Subject}}"
	}
}
