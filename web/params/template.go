package params

import (
	"encoding/json"
	"io"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type Template struct {
	Name     string          `json:"name"`
	Text     string          `json:"text"`
	HTML     string          `json:"html"`
	Subject  string          `json:"subject"`
	Metadata json.RawMessage `json:"metadata"`
}

type TemplateCreateError struct{}

func (err TemplateCreateError) Error() string {
	return "Failed to create Template in the database"
}

func NewTemplate(body io.Reader) (Template, error) {
	var template Template
	err := json.NewDecoder(body).Decode(&template)
	if err != nil {
		return Template{}, ParseError{}
	}
	if template.Metadata == nil {
		template.Metadata = json.RawMessage("{}")
	}

	err = template.Validate()
	if err != nil {
		return Template{}, err
	}

	template.setDefaults()

	return template, nil
}

func (template Template) Validate() error {
	if template.Name == "" {
		return ValidationError([]string{"Request is missing the required field: name"})
	}

	if template.HTML == "" {
		return ValidationError([]string{"Request is missing the required field: html"})
	}

	return nil
}

func (template Template) ToModel() models.Template {
	return models.Template{
		Name:     template.Name,
		Text:     template.Text,
		HTML:     template.HTML,
		Subject:  template.Subject,
		Metadata: string(template.Metadata),
	}
}

func (template *Template) setDefaults() {
	if template.Subject == "" {
		template.Subject = "{{.Subject}}"
	}
}
