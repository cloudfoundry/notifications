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

type roundtripRecorder interface {
	Record(key string, req *http.Request, resp *http.Response) error
}

type Config struct {
	Host              string
	Trace             bool
	Routerless        bool
	RoundTripRecorder roundtripRecorder
}

type Client struct {
	config                  Config
	httpClient              *http.Client
	documentNextRequest     bool
	documentationRequestKey string
}

func NewClient(config Config) *Client {
	return &Client{
		config:     config,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) Do(method string, path string, payload map[string]interface{}, token string) (int, map[string]interface{}, error) {
	var results map[string]interface{}
	responseCode, err := c.DoTyped(method, path, payload, token, &results)
	return responseCode, results, err
}

func (c *Client) DoTyped(method string, path string, payload map[string]interface{}, token string, results interface{}) (int, error) {
	var requestBody io.ReadSeeker

	if payload != nil {
		content, err := json.Marshal(payload)
		if err != nil {
			return 0, err
		}

		requestBody = bytes.NewReader(content)
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

func (c *Client) Document(requestKey string) {
	c.documentNextRequest = true
	c.documentationRequestKey = requestKey
}

func (c *Client) makeRequest(method, path string, content io.ReadSeeker, token string) (int, []byte, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return 0, []byte{}, err
	}
	c.printRequest(request)

	request.Header.Set("X-NOTIFICATIONS-VERSION", "2")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
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

	if c.documentNextRequest {
		if content != nil {
			_, err := content.Seek(0, 0)
			if err != nil {
				return 0, []byte{}, err
			}

			rc, ok := content.(io.ReadCloser)
			if !ok && content != nil {
				rc = ioutil.NopCloser(content)
			}
			request.Body = rc
		}

		response.Body = ioutil.NopCloser(bytes.NewReader(responseBody))

		if err := c.config.RoundTripRecorder.Record(c.documentationRequestKey, request, response); err != nil {
			return 0, []byte{}, err
		}

		c.documentNextRequest = false
	}

	return response.StatusCode, responseBody, nil
}

func (c Client) printRequest(request *http.Request) {
	if !c.config.Trace {
		return
	}

	var buffer *bytes.Buffer
	request.Body, buffer = duplicateBuffer(request.Body)

	fmt.Printf("[REQ] %s %s %s\n", request.Method, request.URL.String(), buffer.String())
}

func (c Client) printResponse(response *http.Response) {
	if !c.config.Trace {
		return
	}

	var buffer *bytes.Buffer
	response.Body, buffer = duplicateBuffer(response.Body)

	fmt.Printf("[RES] %s %s\n", response.Status, buffer.String())
}

func duplicateBuffer(originalBody io.Reader) (io.ReadCloser, *bytes.Buffer) {
	buffer := bytes.NewBuffer([]byte{})
	body := bytes.NewBuffer([]byte{})
	if originalBody != nil {
		_, err := io.Copy(io.MultiWriter(buffer, body), originalBody)
		if err != nil {
			panic(err)
		}
	}
	return ioutil.NopCloser(body), buffer
}
