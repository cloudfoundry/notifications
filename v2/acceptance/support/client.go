package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Config struct {
	Host       string
	Trace      bool
	Routerless bool
}

type Client struct {
	config     Config
	httpClient *http.Client
}

func NewClient(config Config) *Client {
	return &Client{
		config:     config,
		httpClient: http.DefaultClient,
	}
}
func (c Client) DoTyped(method string, path string, payload map[string]interface{}, token string, results interface{}) (int, error) {
	var requestBody io.Reader

	if payload != nil {
		content, err := json.Marshal(payload)
		if err != nil {
			return 0, err
		}

		requestBody = bytes.NewBuffer(content)
	}

	responseCode, responseBody, err := c.makeRequest(method, c.config.Host+path, requestBody, token)
	if err != nil {
		return 0, err
	}

	if strings.Contains(string(responseBody), "{") {
		err = json.Unmarshal(responseBody, &results)
		if err != nil {
			return 0, err
		}
	}

	return responseCode, err
}

func (c Client) Do(method string, path string, payload map[string]interface{}, token string) (int, map[string]interface{}, error) {
	var requestBody io.Reader

	if payload != nil {
		content, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, err
		}

		requestBody = bytes.NewBuffer(content)
	}

	responseCode, responseBody, err := c.makeRequest(method, c.config.Host+path, requestBody, token)
	if err != nil {
		return 0, nil, err
	}

	var body map[string]interface{}

	if strings.Contains(string(responseBody), "{") {
		err = json.Unmarshal(responseBody, &body)
		if err != nil {
			return 0, nil, err
		}
	}

	return responseCode, body, err
}

func (c Client) makeRequest(method, path string, content io.Reader, token string) (int, []byte, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return 0, []byte{}, err
	}
	c.printRequest(request)

	request.Header.Set("X-NOTIFICATIONS-VERSION", "2")
	request.Header.Set("Authorization", "Bearer "+token)
	if c.config.Routerless {
		request.Header.Set("X-Vcap-Request-Id", "some-totally-fake-vcap-request-id")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return 0, []byte{}, err
	}
	defer response.Body.Close()

	c.printResponse(response)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return response.StatusCode, responseBody, nil
}

func (c Client) printRequest(request *http.Request) {
	if !c.config.Trace {
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	body := bytes.NewBuffer([]byte{})
	if request.Body != nil {
		_, err := io.Copy(io.MultiWriter(buffer, body), request.Body)
		if err != nil {
			panic(err)
		}
	}

	request.Body = ioutil.NopCloser(body)

	fmt.Printf("[REQ] %s %s %s\n", request.Method, request.URL.String(), buffer.String())
}

func (c Client) printResponse(response *http.Response) {
	if !c.config.Trace {
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	body := bytes.NewBuffer([]byte{})
	if response.Body != nil {
		_, err := io.Copy(io.MultiWriter(buffer, body), response.Body)
		if err != nil {
			panic(err)
		}
	}

	response.Body = ioutil.NopCloser(body)

	fmt.Printf("[RES] %s %s\n", response.Status, buffer.String())
}
