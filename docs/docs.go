package docs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/pivotal-cf-experimental/warrant"
)

type MethodRequest struct {
	Headers map[string][]string
	Body    string
}

type MethodResponse struct {
	Headers map[string][]string
	Code    int
	Body    string
}

type MethodEntry struct {
	Verb        string
	Description string
	Request     MethodRequest
	Responses   []MethodResponse
}

type ResourceEntry struct {
	ListResourceName  string
	ItemResourceName  string
	ListMethodEntries []MethodEntry
	ItemMethodEntries []MethodEntry
}

type DocGenerator struct {
	RequestInspector         requestInspector
	Resources                map[string]ResourceEntry
	SampleAuthorizationToken string
	DecodedToken             string
}

type requestInspector interface {
	GetResourceInfo(request *http.Request) ResourceInfo
}

func NewDocGenerator(requestInspector requestInspector) *DocGenerator {
	return &DocGenerator{
		Resources:        map[string]ResourceEntry{},
		RequestInspector: requestInspector,
	}
}

func (g *DocGenerator) Add(request *http.Request, response *http.Response) error {
	var resourceEntry ResourceEntry
	var requestBody []byte

	resourceInfo := g.RequestInspector.GetResourceInfo(request)

	if retrievedResource, ok := g.Resources[resourceInfo.ResourceType]; ok {
		resourceEntry = retrievedResource
	}

	resourceEntry.ListResourceName = resourceInfo.ListName
	resourceEntry.ItemResourceName = resourceInfo.ItemName

	if request.Body != nil {
		var err error
		requestBody, err = ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	methodEntry := MethodEntry{
		Verb: request.Method,
		Request: MethodRequest{
			Headers: request.Header,
			Body:    string(requestBody),
		},
		Responses: []MethodResponse{
			{
				Code:    response.StatusCode,
				Headers: response.Header,
				Body:    string(responseBody),
			},
		},
	}

	if resourceInfo.IsItem {
		resourceEntry.ItemMethodEntries = append(resourceEntry.ItemMethodEntries, methodEntry)
	} else {
		resourceEntry.ListMethodEntries = append(resourceEntry.ListMethodEntries, methodEntry)
	}

	g.Resources[resourceInfo.ResourceType] = resourceEntry

	g.SampleAuthorizationToken = resourceInfo.AuthorizationToken

	return nil
}

func (g *DocGenerator) GenerateBlueprint(outputFilePath string) error {
	g.DecodedToken = decodeSampleAuthorizationToken(g.SampleAuthorizationToken)

	if outputFilePath == "" {
		return nil
	}

	tmpl, err := template.ParseFiles("../../docs/template.tmpl")
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(outFile, g)
	if err != nil {
		panic(err)
	}

	outFile.Close()
	return nil
}

func decodeSampleAuthorizationToken(userToken string) string {
	w := warrant.New(warrant.Config{})

	decodedToken, err := w.Tokens.Decode(userToken)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("ClientID: %s\nScopes: %+v\n", decodedToken.ClientID, decodedToken.Scopes)
}
