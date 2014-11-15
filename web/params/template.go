package params

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type Template struct {
	Name string `json:"name"`
	Text string `json:"text"`
	HTML string `json:"html"`
}

type TemplateUpdateError struct{}

func (err TemplateUpdateError) Error() string {
	return "failed to update Template in the database"
}

func NewTemplate(templateName string, body io.Reader) (Template, error) {
	var template Template

	jsonBody, err := ioutil.ReadAll(body)
	if err != nil {
		return Template{}, ParseError{}
	}

	err = json.Unmarshal(jsonBody, &template)
	if err != nil {
		return template, ParseError{}
	}

	err = containsArguments(string(jsonBody))
	if err != nil {
		return Template{}, err
	}

	template.Name = templateName

	return template, nil
}

func containsArguments(jsonBody string) error {
	if !strings.Contains(jsonBody, `"html":`) || !strings.Contains(jsonBody, `"text":`) {
		return ValidationError([]string{"Request is missing a required field"})
	}
	return nil
}

func (template *Template) Validate() error {
	invalidSuffix := true
	name := template.Name

	for _, validEnding := range models.TemplateNames {
		if strings.HasSuffix(name, validEnding) {
			invalidSuffix = false
		}
	}

	if invalidSuffix {
		return ValidationError([]string{fmt.Sprintf("Template has invalid suffix, must end with one of %v\n", models.TemplateNames)})
	}

	return template.validateFormat(name)
}

func (template *Template) validateFormat(name string) error {
	nameParts := strings.Split(name, ".")
	if len(nameParts) == 4 && nameParts[2] != "subject" {
		return ValidationError([]string{"Template name has an invalid format, too many periods."})
	}

	if len(nameParts) > 5 {
		return ValidationError([]string{"Template name has an invalid format, too many periods."})
	}

	if template.hasInvalidCharacters(name) {
		return ValidationError([]string{"Template name has an invalid format, only .-_ and alphanumeric characters are valid."})
	}

	return nil
}

func (template *Template) hasInvalidCharacters(name string) bool {
	replacer := strings.NewReplacer("-", "", ".", "")
	strippedName := replacer.Replace(name)

	regex := regexp.MustCompile(`[\W]`)
	return regex.Match([]byte(strippedName))
}

func (t *Template) ToModel() models.Template {
	template := models.Template{
		Name:       t.Name,
		Text:       t.Text,
		HTML:       t.HTML,
		Overridden: true,
	}
	return template
}
