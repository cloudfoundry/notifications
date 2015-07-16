package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Host       string
	Trace      bool
	Routerless bool
}

type Client struct {
	config     Config
	httpClient *http.Client
	API        APIService
	Senders    SendersService
}

func NewClient(config Config) *Client {
	return &Client{
		config:     config,
		httpClient: http.DefaultClient,
		API:        NewAPIService(config),
		Senders:    NewSendersService(config),
	}
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

type APIService struct {
	config Config
}

func NewAPIService(config Config) APIService {
	return APIService{
		config: config,
	}
}

func (s APIService) Version() (int, error) {
	status, body, err := NewClient(s.config).makeRequest("GET", s.config.Host+"/info", nil, "")
	if err != nil {
		return 0, err
	}

	if status != http.StatusOK {
		return 0, fmt.Errorf("request failed: %d %s", status, body)
	}

	var info struct {
		Version int `json:"version"`
	}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return 0, err
	}

	return info.Version, nil
}
