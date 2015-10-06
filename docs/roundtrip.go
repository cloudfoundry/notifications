package docs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/pivotal-cf-experimental/warrant"
)

type RoundTrip struct {
	Request  *http.Request
	Response *http.Response
}

func (r RoundTrip) Method() string {
	return r.Request.Method
}

func (r RoundTrip) Path() string {
	return guidRegexp.ReplaceAllLiteralString(r.Request.URL.Path, "{id}")
}

func (r RoundTrip) RequiredScopes() string {
	authHeader := r.Request.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ""
	}

	tokensService := warrant.NewTokensService(warrant.Config{})

	token, err := tokensService.Decode(parts[1])
	if err != nil {
		return ""
	}

	return strings.Join(token.Scopes, " ")
}

func (r RoundTrip) RequestHeaders() []string {
	var headers []string

	for key, values := range r.Request.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
	}

	sort.Strings(headers)

	return headers
}

func (r RoundTrip) RequestBody() string {
	if r.Request.Body == nil {
		return ""
	}

	var input map[string]interface{}
	err := json.NewDecoder(r.Request.Body).Decode(&input)
	if err != nil {
		return ""
	}

	output, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return ""
	}

	return string(output)
}

func (r RoundTrip) ResponseStatus() string {
	return r.Response.Status
}

func (r RoundTrip) ResponseHeaders() []string {
	var headers []string

	for key, values := range r.Response.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
	}

	sort.Strings(headers)

	return headers
}

func (r RoundTrip) ResponseBody() string {
	if r.Response.Body == nil {
		return ""
	}

	var input map[string]interface{}
	err := json.NewDecoder(r.Response.Body).Decode(&input)
	if err != nil {
		return ""
	}

	output, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return ""
	}

	return string(output)
}
