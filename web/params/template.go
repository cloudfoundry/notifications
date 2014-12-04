package params

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type Template struct {
	Name    string `json:"name"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
	Subject string `json:"subject"`
}

type TemplateUpdateError struct{}

func (err TemplateUpdateError) Error() string {
	return "failed to update Template in the database"
}

type TemplateCreateError struct{}

func (err TemplateCreateError) Error() string {
	return "Failed to create Template in the database"
}

func NewTemplate(templateID string, body io.Reader) (Template, error) {
	var template Template

	jsonBody, err := ioutil.ReadAll(body)
	if err != nil {
		return Template{}, ParseError{}
	}

	err = json.Unmarshal(jsonBody, &template)
	if err != nil {
		return template, ParseError{}
	}

	if template.Subject == "" {
		template.Subject = "{{.Subject}}"
	}

	err = containsRequiredArguments(string(jsonBody))
	if err != nil {
		return Template{}, err
	}

	return template, nil
}

func (t *Template) ToModel() models.Template {
	template := models.Template{
		Name:    t.Name,
		Text:    t.Text,
		HTML:    t.HTML,
		Subject: t.Subject,
	}
	return template
}

func containsRequiredArguments(jsonBody string) error {
	var validationErrors []string
	if !strings.Contains(jsonBody, `"name":`) {
		validationErrors = append(validationErrors, "Request is missing the required name field")
	}

	if !strings.Contains(jsonBody, `"html":`) {
		validationErrors = append(validationErrors, "Request is missing the required html field")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return ValidationError(validationErrors)
}

//func (template *Template) Validate() error {
//return template.validateFormat(template.ID)
//}

//func (template *Template) validateFormat(id string) error {
//if template.hasInvalidCharacters(id) {
//return ValidationError([]string{"Template id has an invalid format, only .-_ and alphanumeric characters are valid."})
//}

//return nil
//}

//func (template *Template) hasInvalidCharacters(id string) bool {
//replacer := strings.NewReplacer("-", "", ".", "")
//strippedID := replacer.Replace(id)

//regex := regexp.MustCompile(`[\W]`)
//return regex.Match([]byte(strippedID))
//}
