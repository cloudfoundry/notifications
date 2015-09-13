package valiant

import (
	"fmt"
	"strings"
)

func reportExtraKeys(typed interface{}, untyped interface{}) error {
	var expectedTags, actualTags []string
	_, parentType := reflectedParent(typed)
	parentUntyped := untyped.(map[string]interface{})

	for i := 0; i < parentType.NumField(); i++ {
		expectedTags = append(expectedTags, jsonName(parentType.Field(i)))
	}

	for key := range parentUntyped {
		actualTags = append(actualTags, key)
	}

	for _, tag := range actualTags {
		if !contains(expectedTags, tag) {
			return ExtraFieldError{ErrorMessage: fmt.Sprintf("Extra field %q is not valid", tag)}
		}
	}

	return nil
}

func contains(values []string, element string) bool {
	for _, value := range values {
		if strings.ToLower(value) == strings.ToLower(element) {
			return true
		}
	}

	return false
}
