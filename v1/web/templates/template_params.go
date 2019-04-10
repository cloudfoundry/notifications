package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/cloudfoundry-incubator/notifications/valiant"
)

type TemplateParams struct {
	Name     string          `json:"name" validate-required:"true"`
	Text     string          `json:"text"`
	HTML     string          `json:"html" validate-required:"true"`
	Subject  string          `json:"subject"`
	Metadata json.RawMessage `json:"metadata"`
}

func NewTemplateParams(body io.ReadCloser) (TemplateParams, error) {
	defer body.Close()

	var template TemplateParams
	validator := valiant.NewValidator(body)

	err := validator.Validate(&template)
	if err != nil {
		switch err.(type) {
		case valiant.RequiredFieldError:
			return template, webutil.ValidationError{Err: err}
		default:
			return template, webutil.ParseError{}
		}
	}

	if template.Metadata == nil {
		template.Metadata = json.RawMessage("{}")
	}

	err = template.validateSyntax()
	if err != nil {
		return TemplateParams{}, err
	}

	template.setDefaults()

	return template, nil
}

func (t TemplateParams) validateSyntax() error {
	toValidate := map[string]string{
		"Subject": t.Subject,
		"Text":    t.Text,
		"HTML":    t.HTML,
	}

	for field, contents := range toValidate {
		_, err := template.New("test").Parse(contents)
		if err != nil {
			return webutil.ValidationError{Err: fmt.Errorf("%s syntax is malformed please check your braces", field)}
		}
	}

	return nil
}

func (t TemplateParams) ToModel() models.Template {
	return models.Template{
		Name:     t.Name,
		Text:     t.Text,
		HTML:     t.HTML,
		Subject:  t.Subject,
		Metadata: string(t.Metadata),
	}
}

func (t *TemplateParams) setDefaults() {
	if t.Subject == "" {
		t.Subject = "{{.Subject}}"
	}
}
