package docs

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"text/template"
)

const preamble = `<!---
This document is automatically generated.
DO NOT EDIT THIS BY HAND.
Run the acceptance suite to re-generate the documentation.
-->

`

var guidRegexp = regexp.MustCompile(`[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}`)

type TemplateContext struct {
	Resources []TemplateResource
}

type TemplateResource struct {
	Name        string
	Description string
	Endpoints   []TemplateEndpoint
}

type TemplateEndpoint struct {
	Key             string
	Description     string
	Method          string
	Path            string
	RequiredScopes  string
	RequestHeaders  []string
	RequestBody     string
	ResponseStatus  string
	ResponseHeaders []string
	ResponseBody    string
}

func BuildTemplateContext(resources []Resource, roundtrips map[string]RoundTrip) (TemplateContext, error) {
	var context TemplateContext

	for _, resource := range resources {
		templateResource := TemplateResource{
			Name:        resource.Name,
			Description: resource.Description,
		}

		for _, endpoint := range resource.Endpoints {
			roundtrip, ok := roundtrips[endpoint.Key]
			if !ok {
				return TemplateContext{}, fmt.Errorf("missing roundtrip %q", endpoint.Key)
			}

			templateResource.Endpoints = append(templateResource.Endpoints, TemplateEndpoint{
				Key:             endpoint.Key,
				Description:     endpoint.Description,
				Method:          roundtrip.Method(),
				Path:            roundtrip.Path(),
				RequiredScopes:  roundtrip.RequiredScopes(),
				RequestHeaders:  roundtrip.RequestHeaders(),
				RequestBody:     roundtrip.RequestBody(),
				ResponseStatus:  roundtrip.ResponseStatus(),
				ResponseHeaders: roundtrip.ResponseHeaders(),
				ResponseBody:    roundtrip.ResponseBody(),
			})

			delete(roundtrips, endpoint.Key)
		}

		context.Resources = append(context.Resources, templateResource)
	}

	if len(roundtrips) > 0 {
		var keys []string
		for key, _ := range roundtrips {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		return TemplateContext{}, fmt.Errorf("unused roundtrips %v", keys)
	}

	return context, nil
}

func GenerateMarkdown(context TemplateContext) (string, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/docs/template.tmpl", os.Getenv("ROOT_PATH")))
	if err != nil {
		panic(err)
	}

	output := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(output, context)
	if err != nil {
		panic(err)
	}

	return preamble + output.String(), nil
}
