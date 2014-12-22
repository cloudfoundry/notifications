package params

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type DeprecatedTemplate struct {
	Name string `json:"name"`
	Text string `json:"text"`
	HTML string `json:"html"`
}

func NewDeprecatedTemplate(templateName string, body io.Reader) (DeprecatedTemplate, error) {
	var template DeprecatedTemplate

	decodeReader := bytes.NewBuffer([]byte{})
	stringReader := bytes.NewBuffer([]byte{})
	io.Copy(io.MultiWriter(decodeReader, stringReader), body)

	err := json.NewDecoder(decodeReader).Decode(&template)
	if err != nil {
		return DeprecatedTemplate{}, ParseError{}
	}

	err = deprecatedContainsArguments(stringReader.String())
	if err != nil {
		return DeprecatedTemplate{}, err
	}

	template.Name = templateName

	return template, nil
}

func deprecatedContainsArguments(jsonBody string) error {
	if !strings.Contains(jsonBody, `"html":`) || !strings.Contains(jsonBody, `"text":`) {
		return ValidationError([]string{"Request is missing a required field"})
	}
	return nil
}

func (template *DeprecatedTemplate) Validate() error {
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

func (template *DeprecatedTemplate) validateFormat(name string) error {
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

func (template *DeprecatedTemplate) hasInvalidCharacters(name string) bool {
	replacer := strings.NewReplacer("-", "", ".", "")
	strippedName := replacer.Replace(name)

	regex := regexp.MustCompile(`[\W]`)
	return regex.Match([]byte(strippedName))
}

func (t *DeprecatedTemplate) ToModel() models.Template {
	template := models.Template{
		Name: t.Name,
		Text: t.Text,
		HTML: t.HTML,
	}
	return template
}
